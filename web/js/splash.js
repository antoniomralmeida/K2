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


(function () {
  $(window).load(function () {
      setTimeout(function () {
         $('.loader').fadeOut();
          start();
      }, 200);
  });
  $(document).ready(function () {
      $('.splash').slideToggle(3000);
      $('.frame1').delay(3200).show('slide', { direction: 'right' }, 1000);
      $('.frame2').delay(3300).show('slide', { direction: 'left' }, 1000);
      $('.frame3').delay(3800).slideToggle(1000);
      $('.frame3').css('line-height', '400px');
      $('.frame2').delay(4000).hide('slide', { direction: 'left' }, 1000);
      $('.frame1').delay(4200).hide('slide', { direction: 'right' }, 1000);
      $('.splash').delay(4500).slideToggle(2000);
      $('.frame3').delay(4800).transition({
          scale: [
              1.2,
              1.2
          ],
          duration: 1000
      });
      $('.splash').delay(500).animate({ width: '50%' }, 1000);
  });
}.call(this));