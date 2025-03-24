document.addEventListener('DOMContentLoaded', function() {
    console.log("DOMContentLoaded fired")
    init();
});

let errorMsg = '';

function init() {

    loadPage();

    if (userLoggedIn) {
        enableCommentingAndLiking();
    } else {
        disableCommentingAndLiking();
    }
    
    document.getElementById("login-button").addEventListener("click", function () {
        history.pushState({}, '', '/login');  
        loadPage(); 
    });

      // Handle logout
      document.getElementById('logout-button').addEventListener('click', () => {
        userLoggedIn = false; 
        loadPage('login'); 
        disableCommentingAndLiking(); 
    });

    // Handle navigation events (e.g., clicking on links or buttons)
    document.addEventListener("click", async (event) => {
        const postLink = event.target.closest(".post-title a");
        if (!postLink) return;

        // Get the post ID from the dataset and load the page
        const postID = postLink.dataset.postId;

        history.pushState({}, "", `/post/${postID}`);
        loadPage();
        event.preventDefault();

    });

    // Handle back/forward navigation
    window.addEventListener('popstate', () => {
        loadPage(); // Re-run init to reload correct page based on new URL
    });

    

    // document.getElementById('loginLink').addEventListener('click', () => loadPage('login'));
    // document.getElementById('registerLink').addEventListener('click', () => loadPage('register'));
    // document.getElementById('feedLink').addEventListener('click', () => loadPage('feed'));
    // document.getElementById('createPostLink').addEventListener('click', () => loadPage('createPost'));
    // document.getElementById('chatLink').addEventListener('click', () => loadPage('chat'));
}

// toggle which page is shown
function loadPage() {
    hideAllPages();

    const path = window.location.pathname; // Get the URL
    const segments = path.split('/').filter(Boolean); // Remove empty segments
    let page;
    console.log(segments)

    if (segments.length === 0) {
        page = 'feed'
        loadFeedPage();
    } else  if (segments[0] === 'post') {
        page = 'post-details'
        loadPostPage();
    } else if (segments[0] === 'login') {
        page = 'login-page'
    }

    showPage(page)

}

function showPage(pageId) {
    if (pageId === null) return;
    // Show the page and hide others
    document.getElementById(pageId).style.display = 'block';
}

function hideAllPages() {
    const pages = document.querySelectorAll('.page');
    pages.forEach(page => page.style.display = 'none');
}

function showError() {
    const errorContainer = document.querySelector("#error-message");
    const errorText = document.querySelector("#error-text");
    const backButton = document.querySelector("#error-back-btn");

    if (!errorContainer || !errorText || !backButton) return;

    errorText.textContent = errorMsg;

    backButton.addEventListener("click", () => {
        loadPage('feed');
    });
}
