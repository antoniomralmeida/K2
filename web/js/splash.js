let cookieModal = document.querySelector(".cookie-consent-modal")
let cancelCookieBtn = document.querySelector(".btn.cancel")
let acceptCookieBtn = document.querySelector(".btn.accept")
let url = window.location.href +'home';

cancelCookieBtn.addEventListener("click", function () {
  cookieModal.classList.remove("active")
})
acceptCookieBtn.addEventListener("click", function () {
  cookieModal.classList.remove("active")
  localStorage.setItem("cookieAccepted", "yes")
  window.location.href = url;
})

setTimeout(function () {
  let cookieAccepted = localStorage.getItem("cookieAccepted")
  if (cookieAccepted != "yes") {
    cookieModal.classList.add("active");
  } else {
    window.location.href = url;
  }
}, 2000)