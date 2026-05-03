function validateSearch(form) {
  const query = form.q.value.trim();
  if (query === "") {
    alert("Please enter a search query.");
    return false;
  }
  return true;
}