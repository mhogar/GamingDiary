var toggles = document.getElementsByClassName("list-toggle");

for (var i = 0; i < toggles.length; i++) {
    toggles[i].addEventListener("click", function() {
        var nodes = this.parentElement.querySelectorAll(".collapsable");

        for (var j = 0; j < nodes.length; j++) {
            nodes[j].classList.toggle("show");
        }
    });
} 
