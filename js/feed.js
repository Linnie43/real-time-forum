function loadFeedPage() {
    fetch('/api/feed', {
        method: 'GET',
        //credentials: 'same-origin'  // Ensure the session cookie is sent along with the request
    })
    .then(response => {
        if (!response.ok) {
            return response.json().then(err => { throw new Error(err.error || "Unknown error"); });
        }
        return response.json();  // Parse the JSON response
    })
    .then(posts => {
        renderPosts(posts);  // Pass posts to the render function
    })
    .catch(error => {
        console.error('Error fetching feed:', error);
        errorMsg = error.message;
        loadPage('error');
    });
}

// Function to render posts on the page
function renderPosts(posts) {
    const feedContainer = document.getElementById('posts-list');
    feedContainer.innerHTML = ''; // Clear the current content

    if (posts.length === 0) {
        feedContainer.innerHTML = '<p>No posts available</p>';
        return;
    }

    // Loop through posts and create HTML for each post
    posts.forEach(post => {
        const postElement = document.createElement('div');
        postElement.classList.add('post');

        // Customize the post element with post data
        postElement.innerHTML = `
            <h3 class="post-title">
                <a href="#" data-post-id="${post.post_id}">${post.post_title}</a>
            </h3>
            <div class="category-container">
                ${post.categories.map(category => `<p class="category-tags">${category}</p>`).join('')}
            </div>
            <div class="post-info-home">
                <span class="material-symbols-outlined" style="font-size: 24px;">filter_vintage</span>
                <span class="username">${post.username}</span>
            </div>
            <p class="post-content">${post.post_content.substring(0, 150)}...</p>
            <div class="icons-container">
                <div class="reaction-buttons">
                    <span class="comment-icon"><span class="material-symbols-outlined">chat</span>${(post.comments ?? []).length}</span>
                    <span class="material-symbols-outlined">thumb_up</span>
                    <span class="reaction-count">${post.likes}</span>
                    <span class="material-symbols-outlined">thumb_down</span>
                    <span class="reaction-count">${post.dislikes}</span>
                </div>
            </div>
        `;

        // Append the post to the feed container
        feedContainer.appendChild(postElement);
    });
}
