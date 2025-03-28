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

  // Login the user using the user data
  await postData("/login", user).then(() => {
    // After logging in, assign current user to user and redirect to the home page
    checkSession();
  });
}

async function register(user) {
  if (user === undefined) {
    throw new Error("Invalid user");
  }

  // Register the user using the user data
  await postData("/register", user).then(() => {
    // After registering, log in the user and redirect to the home page
    login(user);
  });
}

async function logout() {
  if (socket) {
    socket.close();
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
        // If the user is on the chat page, add the message to the chat
        if (chat) {
          chat.receiveMessage(newMsg);
        }
        // Otherwise, add a notification glow to the user
        else {
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
        } else {
          USER_LIST.addTypingIndicator(newMsg.sender_id);
        }
      }
    };
  }
}

document.querySelector(".logo").addEventListener("click", () => {
  checkSession().then(() => {
    if (user) {
      localStorage.removeItem("post");
      document.querySelector(".logo").innerHTML = "FORUM";
      document.querySelector("post-board").render();
    }
  });
});

document.addEventListener("DOMContentLoaded", () => {
  document.querySelector(".nav-right").appendChild(USERNAME_ELEMENT);
  document.querySelector(".nav-right").appendChild(AUTH_BUTTON);
  // Check if the user is logged in
  checkSession();
  document.querySelector(".container").appendChild(page);
});
