// Show/hide sections dynamically
function showLogin() {
    document.getElementById('login-section').style.display = 'block';
    document.getElementById('register-section').style.display = 'none';
    document.getElementById('forum-section').style.display = 'none';
}

function showRegister() {
    document.getElementById('login-section').style.display = 'none';
    document.getElementById('register-section').style.display = 'block';
    document.getElementById('forum-section').style.display = 'none';
}

function showForum() {
    document.getElementById('login-section').style.display = 'none';
    document.getElementById('register-section').style.display = 'none';
    document.getElementById('forum-section').style.display = 'block';
}

// Fake login function (replace with actual API call)
async function login() {
    const username = document.getElementById('login-username').value;
    const password = document.getElementById('login-password').value;

    const response = await fetch('/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password })
    });

    if (response.ok) {
        showForum();
        fetchPosts();
    } else {
        alert("Invalid login");
    }
}

// Fake registration function (replace with actual API call)
async function register() {
    const nickname = document.getElementById('register-nickname').value;
    const age = document.getElementById('register-age').value;
    const gender = document.getElementById('register-gender').value;
    const firstname = document.getElementById('register-firstname').value;
    const lastname = document.getElementById('register-lastname').value;
    const email = document.getElementById('register-email').value;
    const password = document.getElementById('register-password').value;

    const response = await fetch('/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ nickname, age, gender, firstname, lastname, email, password })
    });

    if (response.ok) {
        alert("Registered! Please log in.");
        showLogin();
    } else {
        alert("Registration failed");
    }
}

// Create post (already implemented earlier)
async function createPost() {
    const title = document.getElementById('post-title').value;
    const content = document.getElementById('post-content').value;
    const category = document.getElementById('post-category').value;

    const response = await fetch('/create-post', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ title, content, category })
    });

    const data = await response.json();
    alert(data.message);
    fetchPosts();
}

// Fetch posts (already implemented earlier)
async function fetchPosts() {
    const response = await fetch('/posts');
    const posts = await response.json();

    const postContainer = document.getElementById('post-feed');
    postContainer.innerHTML = '';

    posts.forEach(post => {
        const postElement = document.createElement('div');
        postElement.classList.add('post');
        postElement.innerHTML = `
            <h3>${post.title}</h3>
            <p>${post.content}</p>
            <small>Category: ${post.category} | ${new Date(post.created_at).toLocaleString()}</small>
        `;
        postContainer.appendChild(postElement);
    });
}
