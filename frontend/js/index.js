// Global variables
const AUTH_BUTTON = new AuthenticationButton();
const USERNAME_ELEMENT = document.createElement("span");
USERNAME_ELEMENT.classList.add("username");

// Local storage page references
const PAGES = {
  AuthenticationPage: new AuthenticationPage("login"),
  HomePage: new HomePage(),
};

let page = PAGES[localStorage.getItem("page")] || PAGES["AuthenticationPage"];
let user = null;
let chat = null;
let socket = null;

// Checks if the user has a session cookie and update the page accordingly
async function checkSession() {
  const response = await fetch("/session", {
    method: "GET",
  });

  const data = await response.json().catch((error) => {
    console.log(error);
  });
  // If the response is successful, the user is logged in
  if (response.status === 200) {
    user = data;
    USERNAME_ELEMENT.innerHTML = user ? user.username : "";

    // If the user is on the login page, redirect to the home page
    if (page instanceof AuthenticationPage) {
      switchPage(new HomePage());
    }

    if (!socket) {
      startWS();
    }
    return true;
  }
  // If the response is unsuccessful, the user is not logged in
  else if (!(page instanceof AuthenticationPage) || response.status === 401) {
    if (socket) {
      socket.close();
      socket = null;
    }
    user = null;
    switchPage(new AuthenticationPage("login"));
    return false;
  }
}

async function switchPage(newPage) {
  // Cancel if new page is same as the current page
  if (newPage.constructor.name === page.constructor.name) {
    return;
  }

  const CONTAINER = document.querySelector(".container");
  CONTAINER.innerHTML = "";

  page = newPage;
  CONTAINER.appendChild(page);
  AUTH_BUTTON.update();
  localStorage.setItem("page", page.constructor.name);
}

async function login(user) {
  if (user === undefined) {
    throw new Error("Invalid user");
  }

  try {
    // Make the login request
    const response = await fetch("/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(user)
    });
    
    // Check if login was successful
    if (response.ok) {
      // After successful login, assign current user and redirect
      checkSession();
    } else {
      // If login failed, display error message in the form
      const errorMsg = "Username/email or password is incorrect.";
      const errorDiv = document.querySelector("#error-message");
      
      // Show the error message and set its text
      errorDiv.textContent = errorMsg;
      errorDiv.classList.remove("hidden");
    }
  } catch (error) {
    console.log("Login error:", error);
    // Also show error message for network errors
    const errorDiv = document.querySelector("#error-message");
    errorDiv.textContent = "An error occurred. Please try again.";
    errorDiv.classList.remove("hidden");
  }
}

async function register(user) {
  if (user === undefined) {
    throw new Error("Invalid user");
  }

  try {
    // Make the registration request
    const response = await fetch("/register", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(user),
    });

    // Handle the response
    if (response.ok) {
      // After registering, log in the user and redirect to the home page
      login(user);
    } else if (response.status === 409) {
      // Handle duplicate username/email error
      const errorDiv = document.querySelector("#error-message");
      errorDiv.textContent = "Username or email already exists. Please try again.";
      errorDiv.classList.remove("hidden");
    } else {
      // Handle other errors
      const errorDiv = document.querySelector("#error-message");
      errorDiv.textContent = "An error occurred during registration. Please try again.";
      errorDiv.classList.remove("hidden");
    }
  } catch (error) {
    console.log("Registration error:", error);
    // Show a generic error message for network errors
    const errorDiv = document.querySelector("#error-message");
    errorDiv.textContent = "An error occurred. Please try again.";
    errorDiv.classList.remove("hidden");
  }
}

async function logout() {
  if (socket) {
    socket.close(1000, "User logged out");
    socket = null;
  }
  USERNAME_ELEMENT.innerHTML = "";
  localStorage.removeItem("post");
  await postData("/logout").then(() => {
    user = null;
    switchPage(new AuthenticationPage("login"));
  });
}

async function startWS() {
  const USER_LIST = document.querySelector("user-list");

  // Check if the browser supports websockets
  if (window["WebSocket"]) {
    socket = new WebSocket(`ws://${document.location.host}/ws`);
    socket.onmessage = async function (event) {
      newMsg = JSON.parse(event.data);

      if (newMsg.msg_type == "msg") {
        // Always process the message if there's an open chat window
        if (chat) {
          chat.receiveMessage(newMsg);
          
          // If this message is from someone other than the current chat partner,
          // also show a notification for that user
          if (newMsg.sender_id !== chat.receiver.id && newMsg.receiver_id === user.id) {
            USER_LIST.addNotification(newMsg.sender_id);
          }
        }
        // No chat window open, just show notification
        else if (newMsg.receiver_id === user.id) {
          USER_LIST.addNotification(newMsg.sender_id);
        }
      } else if (newMsg.msg_type == "online") {
        // Update the online status if the user list is loaded
        if (Object.keys(USER_LIST.users).length > 0) {
          USER_LIST.updateOnlineStatus(newMsg.user_ids);
        } else {
          // Otherwise, wait for the user list to load
          setTimeout(() => {
            USER_LIST.updateOnlineStatus(newMsg.user_ids);
          }, 500);
        }
      } else if (newMsg.msg_type == "typing") {
        // If the user is on the chat page, add the typing indicator
        if (chat && chat.receiver.id === newMsg.sender_id) {
          chat.addTypingIndicator();
        }
      }
    };
  }
}

document.querySelector(".logo").addEventListener("click", () => {
  checkSession().then(() => {
    if (user) {
      localStorage.removeItem("post");
      document.querySelector(".logo").innerHTML = "STUDY HALL";
      document.querySelector("post-board").render();
    }
  });
});

let reloadTimeout;

document.addEventListener("DOMContentLoaded", () => {
  clearTimeout(reloadTimeout); // Clear any previous reload timeout
  reloadTimeout = setTimeout(() => {
    document.querySelector(".nav-right").appendChild(USERNAME_ELEMENT);
    document.querySelector(".nav-right").appendChild(AUTH_BUTTON);
    checkSession();
    document.querySelector(".container").appendChild(page);
  }, 100); // Debounce reloads by 100ms
});
