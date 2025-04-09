class AuthenticationPage extends HTMLElement {
    constructor(type = "login") {
      // default contructor type to login
      super(); // call the parent constructor, HTMLElement. This is required for all custom elements
      this.forms = {
        login: `
        <h1>Welcome</h1>
        <div id="error-message" class="error-message hidden"></div>
        <form id="auth-form" action="/login" method="post">
          <label for="username">Username or Email</label>
          <input type="text" class="input" id="username-input" name="username" required/>
          <label for="password">Password</label>
          <input type="password" class="input" id="password-input" name="password" required minlength="6" maxlength="12"/>
          <br />
          <button class="btn" id="submit-btn" type="submit">SIGN IN</button>
        `,
        register: `
        <h1>Sign Up</h1>
        <div id="error-message" class="error-message hidden"></div>
        <form id="auth-form" action="/register" method="post">
          <label for="first-name">First name </label>
          <input type="text" name="firstname" id="first-name-input" required pattern="\\w{1,16}">
          <label for="last-name">Last name </label>
          <input type="text" name="lastname" id="last-name-input" required pattern="\\w{2,16}">
          <label for="email">Email </label>
          <input type="email" name="email" id="email" required>
          <label for="username">Username</label>
          <input type="text" class="input" id="username-input" name="username" required/>
          <label for="password">Password</label>
          <input type="password" class="input" id="password-input" name="password" required minlength="6" maxlength="12"/>
          <label for="gender">Gender</label>
          <select id="gender-input" name="gender" required>
            <option value="">Select...</option>
            <option value="male">Male</option>
            <option value="female">Female</option>
            <option value="other">Other</option>
          </select>
          <label for="birthdate">Date of Birth</label>
          <input type="date" class="input" id="birthdate-input" name="birthdate" required max="${new Date().toISOString().split('T')[0]}"/>
          <br />
          <button class="btn" id="submit-btn" type="submit">SIGN UP</button>
        </form>
        `,
      };
      this.type = type; // type is set to either login or register
    }

    // METHODS:

    // generateUser is a method to generate a user object from the form data
    generateUser(formData) {
      if (this.type === "login") {
        return {
          username: formData.get("username"),
          password: formData.get("password"),
        };
      } else {
        return {
          id: 0,
          username: formData.get("username"),
          email: formData.get("email"),
          firstname: formData.get("firstname"),
          lastname: formData.get("lastname"),
          gender: formData.get("gender"),
          dob: formData.get("birthdate"),
          password: formData.get("password"),
        };
      }
    }
    // submitData is a method to submit the form data
    submitData() {
      const FORM = this.querySelector("#auth-form");
      const FORM_DATA = new FormData(FORM); // create a new FormData object from the form
  
      const USER = this.generateUser(FORM_DATA);
  
      if (this.type === "login") {
        login(USER);
      } else {
        register(USER);
      }
    }
  
    displayForm() {
      this.innerHTML = this.forms[this._type];
      this.querySelector("#auth-form").onsubmit = (event) => {
        event.preventDefault();
        this.submitData();
      };
    }
  
    get type() { // get type of the form
      return this._type;
    }
  
    set type(type) { // set type of the form
      if (type !== "login" && type !== "register") {
        throw new Error("Invalid type");
      }
  
      this._type = type;
      this.displayForm(); // display the form with the new type
    }
  }
  
  customElements.define("auth-page", AuthenticationPage);
  
  class AuthenticationButton extends HTMLButtonElement {
    constructor() {
      super();
      this.innerHTML = "SIGN UP";
      this.classList.add("btn");
      this.id = "auth-btn";
    }
  
    connectedCallback() {
      this.update();
    }
  
    update() {
      if (page instanceof AuthenticationPage) {
        this.formBehavior(); // if page is an instance of AuthenticationPage, call formBehavior
      } else {
        this.logoutBehavior(); // else call logoutBehavior
      }
    }
    // Toggle between login and register forms on button click
    formBehavior() {
      this.innerText = "SIGN UP";
      this.onclick = () => {
        const AUTH_PAGE = document.querySelector("auth-page");
        const TYPE = AUTH_PAGE.type === "login" ? "register" : "login";
        AUTH_PAGE.type = TYPE;
        this.innerText = TYPE === "login" ? "SIGN UP" : "SIGN IN";
      };
    }
  
    logoutBehavior() {
      this.innerText = "SIGN OUT";
      this.onclick = () => {
        logout();
      };
    }
  }
  
  customElements.define("auth-btn", AuthenticationButton, { extends: "button" });