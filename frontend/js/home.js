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
          <h2>Create Post</h2>
          <post-form></post-form>
          <h2>Posts</h2>
          <div class="post-categories">
              <input type="radio" name="category" id="all" value="all" checked />
              <label for="all">All</label>
              <input type="radio" name="category" id="hobbies" value="hobbies" />
              <label for="hobbies">Hobbies</label>
              <input type="radio" name="category" id="health" value="health" />
              <label for="health">Health</label>
              <input type="radio" name="category" id="tech" value="tech" />
              <label for="tech">Tech</label>
              <input type="radio" name="category" id="music" value="music" />
              <label for="music">Music</label>
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
      this.innerHTML = ``;
      const POST_PAGE = new Post(postData);
      POST_PAGE.classList.add("post-full");
      this.appendChild(POST_PAGE);
  
      (await POST_PAGE.getComments()).forEach((comment) => {
        this.appendChild(comment);
      });
      this.appendChild(new PostForm("comment"));
    }
  }
  
  customElements.define("post-board", PostBoard);
  
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
        <h3>Users</h3>
        <ul id="latest-list">
        </ul>
        <ul id="user-list">
        </ul>
      `;
  
      this.users = await this.getUsers();
    }
  
    async getUsers() {
      const USERS = await getData("/user");
      const USER_OBJECTS = {};
      USERS.forEach((user) => {
        const USER_ELEMENT = new User(user);
        this.querySelector("#user-list").appendChild(USER_ELEMENT);
  
        // Store the user in the collection
        USER_OBJECTS[user.id] = USER_ELEMENT;
      });
      return USER_OBJECTS;
    }
  
    // Add new user to the list alphabetically
    async addUser(user) {
      const USER_ELEMENT = new User(user);
      const USER_LIST = this.querySelector("#user-list");
      const USER_ELEMENTS = USER_LIST.querySelectorAll("user-element");
  
      USER_ELEMENT.online = true;
  
      // Find the correct index to insert the user
      let index = 0;
      for (let i = 0; i < USER_ELEMENTS.length; i++) {
        if (USER_ELEMENTS[i].user.username > user.username) {
          index = i;
          break;
        }
      }
  
      // Insert the user at the correct index
      USER_LIST.insertBefore(USER_ELEMENT, USER_ELEMENTS[index]);
  
      // Store the user in the collection
      this.users[user.id] = USER_ELEMENT;
  
      return USER_ELEMENT;
    }
  
  }
  customElements.define("user-list", UserList);
  