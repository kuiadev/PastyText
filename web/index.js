;(() => {

    let expectingMessage = false
    let identity = ''
    
    const { createApp } = Vue
    createApp({
      data(){
        return {
          pastes: '',
          now: Date.now()
        }
      },
      methods: {
        dial() {
          const conn = new WebSocket(`ws://${location.host}/ws`, `pastytextProtocol`)
      
          conn.addEventListener('close', ev => {
            console.log(`WebSocket Disconnected code: ${ev.code}, reason: ${ev.reason}`, true)
            if (ev.code !== 1001) {
              console.log('Reconnecting in 3s')
              setTimeout(this.dial, 3000)
            }
          })
          conn.addEventListener('open', ev => {
            console.info('websocket connected')
          })
      
          // This is where we handle messages received.
          conn.addEventListener('message', ev => {
            if (typeof ev.data !== 'string') {
              console.error('unexpected message type', typeof ev.data)
              return
            }
    
            this.pastes = JSON.parse(ev.data)
          })
    
          window.addEventListener('paste', (event) => {
            const text = (event.clipboardData || window.clipboardData).getData('text')
            if (text) {
              const msg = {"user": identity,
                "action": "add", 
                "text": text}
              conn.send(JSON.stringify(msg))
            }
          })
        },
        setIdentity() {
          if (localStorage.getItem('identity')) {
              identity = localStorage.getItem('identity')
              console.info("friendly name: ", identity)
              return
          }
  
          fetch('/id')
            .then((response) => response.json())
            .then((data) => {
                console.info("new friendly name: ", data.friendly_name)

                identity = data.friendly_name
                localStorage.setItem("identity", identity)
                localStorage.setItem("ipaddress", data.ipaddress)
            })
            .catch((error) => {
                console.error(error.message)
            })
        },
        isPassword(text) {
          // Check if the text has at least 8 characters
          if (text.length < 8) {
              return false;
          }

          // // Check if the text contains any spaces
          // if (/\s/.test(text)) {
          //   return false;
          // }
      
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
      const pastTimestamp = new Date(date).getTime(); // Example past timestamp
      const differenceInSeconds = Math.floor((pastTimestamp - this.now) / 1000);

      const rtf = new Intl.RelativeTimeFormat('en', { numeric: 'auto', style: 'short' });
      console.log(rtf.format(differenceInSeconds, 'second')); // Example output: "in 27 days"

      let unit = 'second'
      switch (true) {
        case differenceInSeconds > 86400:
        unit = 'day'

      }

      return rtf.format(differenceInSeconds, unit)
    }
    
      
      },
      mounted(){
        this.setIdentity()
        this.dial()
      }

    }).mount('#app')
  })()
  