var toggles = document.getElementsByClassName("list-toggle");

for (var i = 0; i < toggles.length; i++) {
    toggles[i].addEventListener("click", function() {
        var nodes = this.parentElement.querySelectorAll(".collapsable");

        for (var j = 0; j < nodes.length; j++) {
            nodes[j].classList.toggle("show");
        }
    });
} 

document.getElementById("scroll-to-bottom").addEventListener("click", function() {
    window.scrollTo({
        top: document.body.scrollHeight,
        behavior: 'smooth'
    });
});

document.getElementById("scroll-to-top").addEventListener("click", function() {
    window.scrollTo({
        top: 0,
        behavior: 'smooth'
    });
});
