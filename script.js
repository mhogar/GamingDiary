var toggles = document.getElementsByClassName("list-toggle");

for (var i = 0; i < toggles.length; i++) {
  toggles[i].addEventListener("click", function() {
    this.parentElement.querySelector(".content").classList.toggle("active");
  });
} 