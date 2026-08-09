function favouritePost(value, btn) {
  const isFav = btn.dataset.fav === "true";
  const url = isFav ? "/unfavourite" : "/favourite";

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
      btn.textContent = isFav ? "Favourite" : "Unfavourite";
      btn.dataset.fav = (!isFav).toString();
    }
  })
  .catch(err => {
    console.error("Error: ", err);
  })
}