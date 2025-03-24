function loadPostPage() {
    const path = window.location.pathname;
    const segments = path.split('/').filter(Boolean); // Remove empty segments
    let postID;

    if (segments[0] === 'post' && segments[1]) {
        postID = segments[1];
    }

    fetch(`/api/post/${postID}`, {
        method: "GET",
        //credentials: "same-origin"  // Uncomment if authentication is needed
    })
    .then(response => {
        if (!response.ok) {
            return response.json().then(err => { throw new Error(err.error || "Unknown error"); });
        }
        return response.json();  // Parse the JSON response
    })
    .then(post => {
        console.log(post);  // Check the structure of the post object, including Categories
        renderPost(post);  // Pass post data to render function
    })
    .catch(error => {
        console.error("Error fetching post:", error);
        errorMsg = error.message;
        loadPage("error");
    });
}

function renderPost(post) {
    const postContentContainer = document.querySelector("#post-content");
    console.log("Categories:", post.categories);  // Log categories to check if it's an array


    // Render the post content
    postContentContainer.innerHTML = `
        <div class="post-header-like-dislike">
            <h2 class="post-title">${post.post_title}</h2>
            <div class="reaction-buttons">
                <span class="material-symbols-outlined">thumb_up</span> 
                <span class="reaction-count">${post.likes}</span>
                <span class="material-symbols-outlined">thumb_down</span>
                <span class="reaction-count">${post.dislikes}</span>
            </div>
        </div>
        <div class="category-container">
            ${post.categories.map(cat => `<p class="category-tags">${cat}</p>`).join('')}
        </div>
        <div class="post-info">
            <span class="username">${post.username}</span>
            <p>${post.created_at}</p>
        </div>
        <div class="post-card">
            <pre>${post.post_content}</pre>
        </div>
        <h3 class="comment-header">Comments:</h3>
        ${post.comments && post.comments.length ? post.comments.map(comment => `
            <div class="comment">
                <p><strong>${comment.username}</strong>: ${comment.created_at}</p>
                <pre>${comment.comment_content}</pre>
            </div>
        `).join('') : '<p>No comments yet.</p>'}
    `;
}

function enableCommentingAndLiking() {
    const commentButtons = document.querySelectorAll('.comment-button');
    const likeButtons = document.querySelectorAll('.like-button');

    commentButtons.forEach(button => button.disabled = false); 
    likeButtons.forEach(button => button.disabled = false); 
}

function disableCommentingAndLiking() {
    const commentButtons = document.querySelectorAll('.comment-button');
    const likeButtons = document.querySelectorAll('.like-button');

    commentButtons.forEach(button => button.disabled = true); 
    likeButtons.forEach(button => button.disabled = true); 
}
