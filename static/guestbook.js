// JavaScript for guestbook localStorage handling

document.addEventListener("DOMContentLoaded", function() {
  const storedGreetings = JSON.parse(localStorage.getItem("greetings")) || [];
  const entriesContainer = document.getElementById("entries");

  if (storedGreetings.length > 0) {
    storedGreetings.forEach(entry => {
      const entryDiv = document.createElement("div");
      entryDiv.className = "entry";
      entryDiv.innerHTML =
        '<p><b>' + (entry.Author || "An anonymous person") + '</b> wrote:</p>' +
        '<pre>' + entry.Content + '</pre>' +
        '<small>' + entry.Date + '</small>';
      entriesContainer.appendChild(entryDiv);
    });
  } else {
    entriesContainer.innerHTML = "<p>No entries yet. Be the first!</p>";
  }
});

function saveGreeting(author, content) {
  const newGreeting = {
    Author: author || "Anonymous",
    Content: content,
    Date: new Date().toLocaleString()
  };

  const storedGreetings = JSON.parse(localStorage.getItem("greetings")) || [];
  storedGreetings.unshift(newGreeting);
  localStorage.setItem("greetings", JSON.stringify(storedGreetings));
}

function handleFormSubmit(event) {
  event.preventDefault();
  const author = document.getElementById("author").value;
  const content = document.getElementById("content").value;

  if (content.trim() === "") {
    alert("Message cannot be empty.");
    return;
  }

  saveGreeting(author, content);
  event.target.submit();
}