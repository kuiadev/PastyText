;(() => {

    let expectingMessage = false

    let identity = ''
    function setIdentity() {
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
    }

    // dial establishes a WebSocket connection to the server.
    function dial() {
      const conn = new WebSocket(`ws://${location.host}/ws`, `pastytextProtocol`)
  
      conn.addEventListener('close', ev => {
        appendLog(`WebSocket Disconnected code: ${ev.code}, reason: ${ev.reason}`, true)
        if (ev.code !== 1001) {
          appendLog('Reconnecting in 3s', true)
          setTimeout(dial, 3000)
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
        const msg = JSON.parse(ev.data)

        const p = appendLog(ev.data)
        if (expectingMessage) {
          p.scrollIntoView()
          expectingMessage = false
        }
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
    }
    setIdentity()
    dial()

  
    const messageLog = document.getElementById('message-log')
    const publishForm = document.getElementById('publish-form')
    const messageInput = document.getElementById('message-input')
  
    // appendLog appends the passed text to messageLog.
    function appendLog(text, error) {
      const p = document.createElement('p')
      // Adding a timestamp to each message makes the log easier to read.
      p.innerText = `${new Date().toLocaleTimeString()}: ${text}`
      if (error) {
        p.style.color = 'red'
        p.style.fontStyle = 'bold'
      }
      messageLog.append(p)
      return p
    }
    appendLog('Submit a message to get started!')
  
    // onsubmit publishes the message from the user when the form is submitted.
    publishForm.onsubmit = async ev => {
      ev.preventDefault()
  
      const msg = messageInput.value
      if (msg === '') {
        return
      }
      messageInput.value = ''
  
      expectingMessage = true
      try {
        const resp = await fetch('/publish', {
          method: 'POST',
          body: msg,
        })
        if (resp.status !== 202) {
          throw new Error(`Unexpected HTTP Status ${resp.status} ${resp.statusText}`)
        }
      } catch (err) {
        appendLog(`Publish failed: ${err.message}`, true)
      }
    }
  })()
  