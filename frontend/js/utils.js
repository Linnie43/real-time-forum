//POST fetch function
async function postData(url = "", data = {}) {
  const response = await fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  }).catch((error) => console.log(error));
  
  // Return the response object
  return response;
}
  
  //GET fetch function
  async function getData(url = "") {
    const response = await fetch(url, {
      method: "GET",
    });
    return response.json();
  }
  