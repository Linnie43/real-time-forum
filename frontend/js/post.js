class Post extends HTMLElement {
  constructor(postData) {
    super();
    this.classList.add("post");
    this.postData = postData;
  }

  connectedCallback() {
    this.render();
  }

  async render() {
    this.innerHTML = `
      <div class="post-header">
        ${
          this.postData.title === undefined
            ? ""
            : `<h3>${this.postData.title}</h3>`
        }
        <p>${await getData("/user?id=" + this.postData.user_id).then(
          (data) => data.username
        )}</p>
      </div>
      <div class="post-body">
        <p>${this.postData.content}</p>
      </div>
      <div class="post-footer">
        <p>${new Date(this.postData.date).toLocaleString()}</p>
      </div>
    `;
  }

  async getComments() {
    const COMMENTS_DATA = await getData(
      "/comment?param=post_id&data=" + this.postData.id
    );
  
    // If there are no comments, return null
    if (!COMMENTS_DATA || COMMENTS_DATA.length === 0) {
      return null;
    }

    // Create a container for the comments
    const COMMENTS_CONTAINER = document.createElement("div");
    COMMENTS_CONTAINER.classList.add("comments-container");
  
    COMMENTS_DATA.forEach((commentData) => {
      const COMMENT = new Post(commentData);
      COMMENT.classList.remove("post");
      COMMENT.classList.add("post-full", "comment");
      COMMENTS_CONTAINER.appendChild(COMMENT); // Append each comment to the container
    });
  
    return COMMENTS_CONTAINER; // Return the container with all comments
  }
}

customElements.define("post-element", Post);

class PostForm extends HTMLElement {
  constructor(type = "post") {
    super();
    this.type = type;
  }

  connectedCallback() {
    this.render();
  }

  async render() {
    this.innerHTML = `
      <form id="post-form">
      ${
        this.type === "post"
          ? '<input type="text" name="title" placeholder="Title" required />'
          : ""
      }
        <textarea name="content" placeholder="Content" required maxlength="800"></textarea>
        ${
          this.type === "post"
            ? `
            <select class="category-dropdown" name="category" required>
              <option value="" disabled selected>-</option>
              <option value="productivity">Productivity</option>
              <option value=feedback">Feedback</option>
              <option value="help">Help</option>
              <option value="resources">Resources</option>
              <option value="fun">Fun</option>
            </select>`
            : ""
        }
        <button class="btn" type="submit">${
          this.type === "post" ? "PUBLISH" : "COMMENT"
        }</button>
      </form>
    `;

    this.querySelector("#post-form").onsubmit = (event) => {
      event.preventDefault();
      this.submitData();
    };
  }

  async submitData() {
    const FORM = this.querySelector("#post-form");
    const FORM_DATA = new FormData(FORM);

    const POST = {
      category: FORM_DATA.get("category"),
      title: FORM_DATA.get("title"),
      content: FORM_DATA.get("content"),
    };

    if (this.type === "post") {
      await postData("/post", POST);
    } else {
      POST.user_id = user.id;
      POST.post_id = JSON.parse(localStorage.getItem("post")).id;
      await postData("/comment", POST);
    }
    document.querySelector("post-board").render();
  }
}

customElements.define("post-form", PostForm);
