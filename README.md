## 📄 PastyText: A simple text sharing tool

PastyText is an open-source tool designed to make sharing text between devices on the same network easy and efficient. 

### Table of Contents
1. [📥 Installation](#installation)
2. [🚀 Features](#features)
3. [🛠️ Use Cases](#usecases)
4. [❓ FAQs](#faqs)

---
## 📥 Installation <a name="installation"></a>

### Docker Compose

These instructions assume the usage of Caddy, but you could use whatever reverse proxy server that suits your needs.

1. Copy the compose.yaml file from the PastyText source code to a new folder
    
2. In this folder, also create a file named Caddyfile with the following content
    
    ```json
    127.0.0.1:80 {
        reverse_proxy pastytext:8080
    }
    ```
    
    1. If you’re hosting this on the web and want to take advantage of HTTPS, your Caddyfile can instead look like the following (replace “example.com” with your domain or subdomain). [Be sure](https://caddyserver.com/docs/quick-starts/https) to update your domain’s A/AAAA records in your DNS provider to point to your server.
        
        ```json
        example.com {
            reverse_proxy pastytext:8080
        }
        ```
        
3. Run the following command, replace `XXX` with the desired release version (e.g. 1.0.0). `-d` utilizes detached mode to run in the background (you’ll probably want to use this mode on your server).
    
    ```bash
    PT_VERSION=XXX docker compose up -d
    ```
    

### Run from source code

1. Navigate to the folder in which you downloaded the PastyText source code.
2. Run the following command.
    
    ```bash
    go run main.go
    ```
    

<aside>
💡

You may need to update the `dial()` function in the index.js file by replacing `wss://` with `ws://` in order for web sockets to run locally.

</aside>

---

## 🚀 Features <a name="features"></a>

| Feature | Description |
| --- | --- |
| **Easy Text Sharing** | Share text snippets with anyone on the same network through a self-hosted page (e.g., [pastytext.example.com](http://pastytext.example.com/)). |
| **Real-Time Updates** | The shared page updates automatically to show new pastes without needing to refresh. New pastes are marked until the page is refreshed or a newer paste is added. |
| **Device Identification** | Automatically assigns unique names to devices on the network (e.g., tasty-wombat) for easy identification of who shared what. |
| **Individual Snippet Management** | Each pasted snippet can be copied or deleted individually, with timestamps indicating when they were shared. |
| **Self-Hosted** | PastyText can be hosted on your own server, ensuring privacy and control over your data. |
| **Plain Text Format** | Maintains formatting for copy-pasted content. |
| **Open Source** | Licensed under the AGPL-3.0, allowing anyone to contribute and improve the tool. |

---

## 🛠️ Use Cases for PastyText <a name="usecases"></a>

- **Collaborative Work**: Share notes, ideas, or code snippets with colleagues in a shared workspace.
- **Family Communication**: Easily share reminders, grocery lists, or messages among family members on the same network.
- **Event Planning**: Coordinate details for events by sharing links, schedules, or tasks with friends.
- **Learning and Education**: Share study materials, resources, or links to educational content with classmates.
- **Quick Access to Links**: Share YouTube links, articles, or other web resources without needing to send them each device individually.

---

## ❓ Frequently Asked Questions (FAQs) <a name="faqs"></a>

### 🧐 How does PastyText work?

PastyText allows users on the same network to share text snippets on a self-hosted page. Users can paste text, which is then visible to anyone on that same network.

### 🧐 Is PastyText secure?

PastyText is not designed for secure sharing of sensitive information like passwords. It operates in plain text, so ensure your network is secure if you choose to share sensitive data.

### 🧐 Can I use PastyText on any device?

Yes! PastyText is a web-based tool that works in any browser, making it accessible on any device that has access to the web/network.

### 🧐 How long does the shared text persist?

The text persists as long as the SQLite storage remains. If you’re using Docker Compose, updating the service may result in purging of the data. Users can delete individual snippets at any time.

### 🧐 Can I share more than just text?

Currently, PastyText only supports text sharing. However, future updates may include support for additional formats.

---
