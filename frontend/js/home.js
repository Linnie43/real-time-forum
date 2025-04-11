class HomePage extends HTMLElement {
  constructor() {
    super();
  }

  async connectedCallback() {
    this.render();
  }

  async render() {
    const POST_BOARD = new PostBoard();
    const USER_LIST = new UserList();
    this.appendChild(POST_BOARD);
    this.appendChild(USER_LIST);
  }
}

customElements.define("home-page", HomePage);

class PostBoard extends HTMLElement {
  constructor() {
    super();
  }

  async connectedCallback() {
    this.render();
  }

  async render() {
    let lastPost = localStorage.getItem("post") || null;
    if (lastPost) {
      this.renderPost(JSON.parse(lastPost));
    } else {
      this.renderPosts();
    }
  }

  async renderPosts() {
    this.innerHTML = `
    <div class="posts-section">
      <h2>Create Post</h2>
      <post-form></post-form>
      <h2>Posts</h2>
      <div class="post-categories">
        <span>Filter:</span>
          <input type="radio" name="category" id="all" value="all" checked />
          <label for="all">All</label>
          <input type="radio" name="category" id="productivity" value="productivity" />
          <label for="productivity">Productivity</label>
          <input type="radio" name="category" id="feedback" value="feedback" />
          <label for="feedback">Feedback</label>
          <input type="radio" name="category" id="help" value="help" />
          <label for="help">Help</label>
          <input type="radio" name="category" id="resources" value="resources" />
          <label for="resources">Resources</label>
          <input type="radio" name="category" id="fun" value="fun" />
          <label for="fun">Fun</label>
      </div>
      <div class="post-container">
      </div>
    `;
    document.querySelector(".post-categories").onchange = (event) => {
      const POSTS = document.querySelectorAll("post-element");
      POSTS.forEach((post) => {
        if (event.target.value === "all") {
          post.style.display = "block";
        } else if (post.postData.category !== event.target.value) {
          post.style.display = "none";
        } else {
          post.style.display = "block";
        }
      });
    };
    const POSTS = await getData("/post");
    POSTS.forEach((postData) => {
      const POST_ELEMENT = new Post(postData);
      POST_ELEMENT.addEventListener("click", () => {
        localStorage.setItem("post", JSON.stringify(postData));
        document.querySelector(".logo").innerHTML = "BACK";
        this.renderPost(postData);

        scrollTo(0, 0);
      });

      this.querySelector(".post-container").appendChild(POST_ELEMENT);
    });
  }

  async renderPost(postData) {
    this.innerHTML = `
      <div class="posts-section-post">
        <div class="post-full"></div>
        <div class="comments-wrapper">
          <h4>Comments</h4>
        </div>
      </div>
    `;
  
    const postContainer = this.querySelector(".post-full");
    const commentsWrapper = this.querySelector(".comments-wrapper");
  
    // Add the full post
    const POST_PAGE = new Post(postData);
    postContainer.appendChild(POST_PAGE);
  
    // Add the comments if they exist
    const COMMENTS_CONTAINER = await POST_PAGE.getComments();
    if (COMMENTS_CONTAINER) {
      commentsWrapper.appendChild(COMMENTS_CONTAINER);
    } else {
      // Optionally hide the comments-wrapper if there are no comments
      commentsWrapper.style.display = "none";
    }
  
    // Add a form for new comments
    const postForm = new PostForm("comment");
    commentsWrapper.style.display = "block"; // Ensure the wrapper is visible for the form
    commentsWrapper.appendChild(postForm);
  }
}

customElements.define("post-board", PostBoard);

class ChatMessage extends HTMLLIElement {
  constructor(sender, content, time = null, isSentByYou = false) {
    super();
    this.sender = sender;
    this.content = content;
    this.time = time === null ? null : new Date(time).toLocaleString();
    this.isSentByYou = isSentByYou;
    this.classList.add("chat-message");

    // Add a class based on whether the message is sent by you or received
    if (isSentByYou) {
      this.classList.add("sent");
    } else {
      this.classList.add("received");
    }
  }

  async connectedCallback() {
    this.render();
  }

  async render() {
    if (this.isSentByYou) {
      this.innerHTML = `
      <span class="chat-username">${this.time === null ? "" : `<span class="chat-time">${this.time}</span>`}
      ${this.sender.username}
      </span>
      <span class="chat-content">${this.content}</span>
    `;
    } else {
      this.innerHTML = `
      <span class="chat-username">${this.sender.username}
      ${this.time === null ? "" : `<span class="chat-time">${this.time}</span>`}
      </span>
      <span class="chat-content">${this.content}</span>
    `;
    }
  }
}

customElements.define("chat-message", ChatMessage, { extends: "li" });

class ChatWindow extends HTMLElement {
  constructor(receiver) {
    super();
    this.receiver = receiver;
    this.typingTimer = null;
    
    // check if user is online
    const userElement = document.querySelector(`user-element[user-id="${receiver.id}"]`);
    this.receiverIsOnline = userElement ? userElement.hasAttribute("online") : false;
  }

  async connectedCallback() {
    this.render();

    let lastScrollTop = 0;
    let isScrolling;
    let isTyping = false;

    const CHAT_BODY = this.querySelector(".chat-body");
    const CHAT_LIST = CHAT_BODY.querySelector("#chat-list");

    CHAT_BODY.onscroll = () => {
      if (CHAT_BODY.scrollTop < 100 && CHAT_BODY.scrollTop < lastScrollTop) {
        // 100px buffer zone
        if (!isScrolling) {
          isScrolling = true;
          this.loadMessages();

          // Throttle scroll event
          setTimeout(() => {
            isScrolling = false;
          }, 500);
        }
      }

      lastScrollTop = CHAT_BODY.scrollTop;
    };

    await this.loadMessages();

    if (this.receiverIsOnline) {
      this.querySelector("#chat-input").oninput = () => {
        if (!isTyping) {
          isTyping = true;

          let typingMsg = {
            id: 0,
            sender_id: user.id,
            receiver_id: this.receiver.id,
            msg_type: "typing",
          };

          socket.send(JSON.stringify(typingMsg));

          setTimeout(() => {
            isTyping = false;
          }, 250);
        }
      };
    }

    this.querySelector("#chat-input").focus();

    this.querySelector("#close-chat").onclick = () => {
      chat = null;
      this.remove();
    };

    this.querySelector("#chat-send").onclick = () => {
      if (this.receiverIsOnline) {
        this.sendMessage();
      }
    };

    this.querySelector("#chat-input").onkeydown = (event) => {
      if (event.key === "Enter" && this.receiverIsOnline) {
        event.preventDefault();
        this.sendMessage();
      }
    };

    CHAT_BODY.scrollTop = CHAT_LIST.scrollHeight;
  }

  async render() {
    this.innerHTML = `
      <div class="chat-header">
        <h3>${this.receiver.username}</h3>
        <button id="close-chat">✖</button>
      </div>
      <div class="chat-body">
      <ul id ="user-list">
        <ul id="chat-list">
        </ul>
      </div>
      <div class="chat-footer">
        <textarea id="chat-input" rows="1" maxlength="500" 
          placeholder="${this.receiverIsOnline ? '' : 'User offline'}"
          ${this.receiverIsOnline ? '' : 'disabled'}></textarea>
        <button id="chat-send" ${this.receiverIsOnline ? '' : 'disabled'}>SEND</button>
      </div>
      </div>
    `;
  }

  async loadMessages() {
    const CHAT_BODY = this.querySelector(".chat-body");
    const CHAT_LIST = CHAT_BODY.querySelector("#chat-list");
    const CHAT_MESSAGES = await getData(
      `/message?receiver=${this.receiver.id}&offset=${CHAT_LIST.children.length}`
    );

    // Reverse the messages so they are in chronological order
    CHAT_MESSAGES.reverse();
    CHAT_MESSAGES.forEach((message) => {
      const CHAT_MESSAGE = new ChatMessage(
        message.sender_id === user.id ? user : this.receiver,
        message.content,
        message.date,
        message.sender_id === user.id // Check if the message is sent by you
      );
      // Add the message on top of current messages
      CHAT_LIST.prepend(CHAT_MESSAGE);
    });
  }

  async sendMessage() {
    const CHAT_INPUT = this.querySelector("#chat-input");
    const CHAT_BODY = this.querySelector(".chat-body");
    const CHAT_LIST = CHAT_BODY.querySelector("#chat-list");

    if (!CHAT_INPUT.value || !socket) {
      console.log("No message or socket");
      return;
    }

    const RECEIVER_USER_ELEMENT = document.querySelector(
      `user-element[user-id="${this.receiver.id}"]`
    );

    // Add the user to the top of #latest-list locally
    RECEIVER_USER_ELEMENT.remove();
    const latestList = document.querySelector("#latest-list");
    latestList.prepend(RECEIVER_USER_ELEMENT);

    // Ensure #latest-list is visible
    if (latestList.children.length > 0) {
      latestList.style.display = "flex"; // Reset display to flex
    }
 
    const CHAT_MESSAGE = new ChatMessage(user, CHAT_INPUT.value, new Date(), true); // Sent by you
    CHAT_LIST.appendChild(CHAT_MESSAGE);
    CHAT_BODY.scrollTop = CHAT_LIST.scrollHeight;
    CHAT_INPUT.value = "";
  
    let msgData = {
      id: 0,
      sender_id: user.id,
      receiver_id: this.receiver.id,
      content: CHAT_MESSAGE.content,
      date: "",
      msg_type: "msg",
    };
  
    socket.send(JSON.stringify(msgData));
  }

  async receiveMessage(message) {
    // Check if this message belongs to the current conversation
    // A message belongs to this conversation if:
    // 1. The current user is the receiver and the sender is the chat partner, OR
    // 2. The current user is the sender and the receiver is the chat partner
    if (!((message.sender_id === this.receiver.id && message.receiver_id === user.id) || 
          (message.sender_id === user.id && message.receiver_id === this.receiver.id))) {
        return; // Skip messages not part of this conversation
    }
    
    const CHAT_BODY = this.querySelector(".chat-body");
    const CHAT_LIST = CHAT_BODY.querySelector("#chat-list");
    const CHAT_MESSAGE = new ChatMessage(
      message.sender_id === user.id ? user : this.receiver,
      message.content,
      message.date,
      message.sender_id === user.id // Check if the message is sent by you
    );
    CHAT_LIST.appendChild(CHAT_MESSAGE);
    this.typingTimer = clearTimeout(this.typingTimer);
    if (CHAT_LIST.querySelector(".typing"))
      CHAT_LIST.querySelector(".typing").remove();
    CHAT_BODY.scrollTop = CHAT_LIST.scrollHeight;
  }

  async addTypingIndicator() {
    const CHAT_BODY = this.querySelector(".chat-body");
    const CHAT_LIST = CHAT_BODY.querySelector("#chat-list");

    // Try to find the typing message
    let TYPING_MESSAGE = CHAT_LIST.querySelector(".typing");

    // If it doesn't exist, create it
    if (!TYPING_MESSAGE) {
      TYPING_MESSAGE = new ChatMessage(this.receiver, "Typing...");
      TYPING_MESSAGE.classList.add("typing");
      CHAT_LIST.appendChild(TYPING_MESSAGE);
    }

    // If there's already a timer, clear it
    if (this.typingTimer) {
      clearTimeout(this.typingTimer);
    }

    // Set a new timer to remove the typing message after 2 seconds
    this.typingTimer = setTimeout(() => {
      TYPING_MESSAGE.remove();
      TYPING_MESSAGE = null;
    }, 2000);

    // If the user is scrolled to the bottom above the typing message, scroll to the bottom again
    CHAT_BODY.scrollTop = CHAT_BODY.scrollHeight;
  }
}

customElements.define("chat-window", ChatWindow);

class User extends HTMLElement {
  constructor(user) {
    super();
    this.user = user;
    this.online = false;
    this.notification = false;
    this.typing = false;
    this.typingTimer = null;
    this.setAttribute("user-id", user.id);
  }

  async connectedCallback() {
    this.render();

    this.onclick = () => {
      // Remove the chat window if it exists
      if (chat) {
        chat.remove();
        chat = null;
      }

      // Return if the user is clicking on themselves
      if (this.user.id === user.id) {
        return;
      }

      // Remove the notification if it exists
      this.notification = false;

      const CHAT_WINDOW = new ChatWindow(this.user);
      document.querySelector(".container").appendChild(CHAT_WINDOW);
      chat = CHAT_WINDOW;
    };
  }

  async render() {
    this.innerHTML = `
      <li>${this.user.username}</li>
    `;
  }

  set online(value) {
    if (value) {
      this.setAttribute("online", "");
    } else {
      this.removeAttribute("online");
    }
  }

  set notification(value) {
    if (value) {
      this.setAttribute("notification", "");
    } else {
      this.removeAttribute("notification");
    }
  }

  set typing(value) {
    if (value) {
      this.setAttribute("typing", "");
    } else {
      this.removeAttribute("typing");
    }
  }
}

customElements.define("user-element", User);

class UserList extends HTMLElement {
  constructor() {
    super();
    this.users = {};
  }

  async connectedCallback() {
    this.render();
  }

  async render() {
    this.innerHTML = `
    <div class="chat-sidebar">
      <h3>Users</h3>
      <div class="user-container">
        <ul id="latest-list">
        </ul>
        <ul id="user-list">
        </ul>
      </div>
    </div>
    `;

    this.users = await this.getUsers();

    // Hide #latest-list if it has no child elements
    const latestList = this.querySelector("#latest-list");
    if (latestList && latestList.children.length === 0) {
      latestList.style.display = "none";
    }
  }

  async getUsers() {
    // Fetch all users
    const USERS = await getData("/user");
    let CHATS = { user_ids: [] }; // Default empty array
    
    try {
      // Fetch chat history
      CHATS = await getData(`/chat?user_id=${user.id}`);
      // Ensure user_ids is an array
      if (!CHATS.user_ids || !Array.isArray(CHATS.user_ids)) {
        CHATS.user_ids = [];
      }
    } catch (error) {
      console.error("Error fetching chat data:", error);
      // Continue with default empty array
    }

    const USER_OBJECTS = {};

    // Sort users by last message time or alphabetically if no chat history
    USERS.sort((a, b) => {
      const chatA = CHATS.user_ids.includes(a.id) ? CHATS.user_ids.indexOf(a.id) : Infinity;
      const chatB = CHATS.user_ids.includes(b.id) ? CHATS.user_ids.indexOf(b.id) : Infinity;

      if (chatA !== chatB) {
        return chatA - chatB; // Sort by chat history
      }
      return a.username.localeCompare(b.username); // Fallback to alphabetical order
    });

    USERS.forEach((userData) => {
      // Renamed variable to avoid conflict with global user
      const USER_ELEMENT = new User(userData);
      this.querySelector("#user-list").appendChild(USER_ELEMENT);

      // Store the user in the collection
      USER_OBJECTS[userData.id] = USER_ELEMENT;
    });
    return USER_OBJECTS;
  }

  async addNotification(userId) {
    if (!this.users[userId]) return;
  
    this.users[userId].notification = true;
  
    // Add the user to the top of #latest-list
    this.users[userId].remove();
    this.users[userId].typing = false;
    const latestList = this.querySelector("#latest-list");
    latestList.prepend(this.users[userId]);
  
    // Ensure #latest-list is visible
    if (latestList.children.length > 0) {
      latestList.style.display = "flex"; // Reset display to flex
    }
  }

  
  async addUser(user) {
    const USER_ELEMENT = new User(user);
    USER_ELEMENT.online = true;
    
    // Add new users to the alphabetical section by default
    const USER_LIST = this.querySelector("#user-list");
    const USER_ELEMENTS = USER_LIST.querySelectorAll("user-element");
    
    // Find the correct alphabetical position
    let inserted = false;
    for (let i = 0; i < USER_ELEMENTS.length; i++) {
      if (USER_ELEMENTS[i].user.username.localeCompare(user.username) > 0) {
        USER_LIST.insertBefore(USER_ELEMENT, USER_ELEMENTS[i]);
        inserted = true;
        break;
      }
    }
    
    if (!inserted) {
      USER_LIST.appendChild(USER_ELEMENT);
    }
    
    // Store the user in the collection
    this.users[user.id] = USER_ELEMENT;
    return USER_ELEMENT;
  }

  async updateOnlineStatus(newOnlineUserIds) {
    // Remove online status from all users
    Object.values(this.users).forEach((user) => {
      user.online = false;
    });

    // Add online status to new users
    newOnlineUserIds.forEach((userId) => {
      // If the user is not in the list, add them
      if (this.users[userId] === undefined) {
        getData(`/user?id=${userId}`).then((data) => {
          return this.addUser(data);
        });

        return;
      }

      this.users[userId].online = true;
    });
  }
}

customElements.define("user-list", UserList);
