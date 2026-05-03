function favouritePost(value, btn) {
  const isFav = btn.dataset.fav === "true";
  const url = isFav ? "/favourite" : "/unfavourite";

  const params = new URLSearchParams();
  params.append("id", value);

  fetch(url, {
    method: "POST",
    headers: {
      "Content-type": "application/x-www-form-urlencoded"
    },
    body: params.toString()
  })
  .then(res => {
    if (res.ok) {
      btn.textContent = isFav ? "Unfavourite" : "Favourite";
      btn.dataset.fav = (!isFav).toString();
    }
  })
  .catch(err => {
    console.error("Error: ", err);
  })
}