;(() => {
    const { createApp } = Vue
    createApp({
      data(){
        return {
          conn: null,
          identity:'',
          network: '',
          lastPasteTime: 0,
          pastes: '',
          now: Date.now(),
          showNewBanner: false,
          showDeleteBanner: false,
          showCopyBanner: false,
          showDelayBanner: false
        }
      },
      computed:{
        availablePastes(){
          let cleanedPastes = []
          let latestPasteIdx = localStorage.getItem("latestPasteIdx");

          if (this.pastes === null) {
            return [];
          }

          Array.from(this.pastes).forEach(function(element) {
            if (element.Id > latestPasteIdx) {
              element.isNew = true;
            } else {
              element.isNew = false;
            }

            element.isShown = true;
            cleanedPastes.push(element);
          });

          if (cleanedPastes.length > 0) {
            localStorage.setItem("latestPasteIdx", cleanedPastes[0]["Id"]);
          }
          
          return cleanedPastes;
        }
      },
      methods: {
        dial() {
          this.conn = new WebSocket(`wss://${location.host}/ws`, `pastytextProtocol`);
      
          this.conn.addEventListener('close', ev => {
            console.log(`WebSocket Disconnected code: ${ev.code}, reason: ${ev.reason}`, true);
            if (ev.code !== 1001) {
              window.removeEventListener('paste', this.handlePaste)
              console.log('Reconnecting in 3s');
              setTimeout(this.dial, 3000);
            }
          })
          this.conn.addEventListener('open', ev => {
            console.info('websocket connected');
          })
      
          // This is where we handle messages received.
          this.conn.addEventListener('message', ev => {
            if (typeof ev.data !== 'string') {
              console.error('unexpected message type', typeof ev.data);
              return;
            }
    
            this.pastes = JSON.parse(ev.data);
            if (this.pastes === null || this.pastes.length === 0) {
              console.log('no pastes');
              localStorage.removeItem("latestPasteIdx");
              return;
            }

            if (this.lastPasteTime == 0 || ((Date.now() - this.lastPasteTime) / 1000) > 3) {
              if (localStorage.getItem("latestPasteIdx") !== null && this.pastes[0].Id > localStorage.getItem("latestPasteIdx")) {
                this.showNewBanner = true;
                this.showCopyBanner = false;
                this.showDeleteBanner = false;
                this.showDelayBanner = false;
              }
            }
          })
    
          window.addEventListener('paste', this.handlePaste);
        },
        handlePaste(){
          // Prevent pasting if the last paste was less than 3 seconds ago
          if (((Date.now() - this.lastPasteTime) / 1000) < 3) {
            this.showDelayBanner = true;
            return;
          }

          this.showDelayBanner = false;

          let pastedText = '';
          navigator.clipboard
            .readText()
            .then(clipText => {
              pastedText = clipText;

              if (pastedText) {
                const msg = {"user": this.identity,
                  "action": "add", 
                  "text": pastedText};
                this.conn.send(JSON.stringify(msg));
                this.lastPasteTime = Date.now();
              }
            });
        },
        setIdentity() {
          if (localStorage.getItem('ipaddress')) {
            this.network = localStorage.getItem('ipaddress');
          }

          if (localStorage.getItem('identity')) {
              this.identity = localStorage.getItem('identity');
              return;
          }
  
          fetch('/id')
            .then((response) => response.json())
            .then((data) => {
                console.info("new friendly name: ", data.friendly_name);

                this.identity = data.friendly_name;
                this.network = data.ipaddress;
                localStorage.setItem("identity", this.identity);
                localStorage.setItem("ipaddress", data.ipaddress);
            })
            .catch((error) => {
                console.error(error.message);
            })
        },
        deletePaste(pasteID) {
          const msg = {"action": "delete", 
            "id": pasteID};
          this.conn.send(JSON.stringify(msg));

          this.showDeleteBanner = true;
          this.showCopyBanner = false;
          this.showNewBanner = false;
          this.showDelayBanner = false;
        },
        isPassword(text) {
          // Check if the text has at least 8 characters
          if (text.length < 8) {
              return false;
          }

          // // Check if the text contains any spaces
          if (/\s/.test(text)) {
            return false;
          }
      
          // Check if the text contains at least one lowercase letter
          if (!/[a-z]/.test(text)) {
              return false;
          }
      
          // Check if the text contains at least one uppercase letter
          if (!/[A-Z]/.test(text)) {
              return false;
          }
      
          // Check if the text contains at least one digit
          if (!/\d/.test(text)) {
              return false;
          }
      
          // Check if the text contains at least one special character
          if (!/[!@#$%^&*(),.?":{}|<>]/.test(text)) {
              return false;
          }
      
          // If all conditions are met, the text is a valid password
          return true;
      },
      isURL(text) {
        // Regular expression to check if the text matches a URL pattern with or without a protocol
        const regex = /^(https?|ftp):\/\/[^\s/$.?#].[^\s]*$/i;
    
        // Check if the URL has a protocol (http://, https://, ftp://) or no protocol (e.g., www.example.com)
        const withProtocol = regex.test(text);
    
        // Check if the URL is a valid format without a protocol (e.g., www.example.com or example.com)
        const withoutProtocol = /^[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)+([\/?].*)?$/.test(text);
    
        return withProtocol || withoutProtocol;
    },
    prettyDate(date) {
      const pastTimestamp = new Date(date).getTime();
      const differenceInSeconds = Math.floor((this.now - pastTimestamp) / 1000);
      const rtf = new Intl.RelativeTimeFormat('en', { numeric: 'auto' });

      let unit = 'second';
      let diff = 0;
      switch (true) {
        case differenceInSeconds > 18144000:
          unit = 'month';
          diff = Math.floor(differenceInSeconds/30/7/24/60/60);
          break;
        case differenceInSeconds > 604800:
          unit = 'week';
        // "last week" use case
          if (differenceInSeconds < 691200){
            diff = Math.floor(differenceInSeconds/7/24/60/60);
          } else {
            diff = Math.ceil(differenceInSeconds/7/24/60/60);
          }
          break;
        case differenceInSeconds > 86400:
          unit = 'day';
          diff = Math.floor(differenceInSeconds/24/60/60);
          break;
        case differenceInSeconds > 3600:
          unit = 'hour';
          diff = Math.floor(differenceInSeconds/60/60);
          break;
        case differenceInSeconds > 60:
          unit = 'minute';
          diff = Math.floor(differenceInSeconds/60);
          break;
        case differenceInSeconds > 15:
          unit = 'second';
          diff = Math.floor(differenceInSeconds);
          break;
        default:
          return "Just now";
      }

      return rtf.format(diff*-1, unit);
    },
    async copyContent(valID, txt, vele) {
      try {
        await navigator.clipboard.writeText(txt);
        this.$refs['copy_' + valID][0].textContent = "Copied!";
        vele.isShown = false;

        this.showCopyBanner = true;
        this.showDeleteBanner = false;
        this.showNewBanner = false;
        this.showDelayBanner = false;
      } catch (error) {
        console.error(error.message);
      }
    },
    hideNewBanner() {
      this.showNewBanner = false;
    },
    hideCopyBanner() {
      this.showCopyBanner = false;
    },
    hideDeleteBanner() {
      this.showDeleteBanner = false;
    },
    hideDelayBanner() {
      this.showDelayBanner = false;
    }
      },
      mounted(){
        this.setIdentity();
        this.dial();
      }

    }).mount('#app')
  })()
  