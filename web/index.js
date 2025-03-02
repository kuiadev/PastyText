;(() => {
    const { createApp } = Vue
    createApp({
      data(){
        return {
          identity:'',
          pastes: '',
          now: Date.now()
        }
      },
      computed:{
        availablePastes(){
          let cleanedPastes = []
          Array.from(this.pastes).forEach(function(element) {
            element.isShown = true;
            cleanedPastes.push(element);
          });
          return cleanedPastes;
        }
      },
      methods: {
        dial() {
          const conn = new WebSocket(`ws://${location.host}/ws`, `pastytextProtocol`);
      
          conn.addEventListener('close', ev => {
            console.log(`WebSocket Disconnected code: ${ev.code}, reason: ${ev.reason}`, true);
            if (ev.code !== 1001) {
              console.log('Reconnecting in 3s');
              setTimeout(this.dial, 3000);
            }
          })
          conn.addEventListener('open', ev => {
            console.info('websocket connected');
          })
      
          // This is where we handle messages received.
          conn.addEventListener('message', ev => {
            if (typeof ev.data !== 'string') {
              console.error('unexpected message type', typeof ev.data);
              return;
            }
    
            this.pastes = JSON.parse(ev.data);
          })
    
          window.addEventListener('paste', (event) => {
            const text = (event.clipboardData || window.clipboardData).getData('text');
            if (text) {
              const msg = {"user": this.identity,
                "action": "add", 
                "text": text};
              conn.send(JSON.stringify(msg));
            }
          })
        },
        setIdentity() {
          if (localStorage.getItem('identity')) {
              this.identity = localStorage.getItem('identity');
              console.info("friendly name: ", this.identity);
              return;
          }
  
          fetch('/id')
            .then((response) => response.json())
            .then((data) => {
                console.info("new friendly name: ", data.friendly_name);

                this.identity = data.friendly_name
                localStorage.setItem("identity", this.identity);
                localStorage.setItem("ipaddress", data.ipaddress);
            })
            .catch((error) => {
                console.error(error.message);
            })
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
          diff = Math.floor(differenceInSeconds/7/24/60/60);
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
      } catch (error) {
        console.error(error.message);
      }
    }
    
      },
      mounted(){
        this.setIdentity();
        this.dial();
      }

    }).mount('#app')
  })()
  