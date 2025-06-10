

const params = new URLSearchParams(window.location.search);
if (params.get("submitted") === "true") {
    document.addEventListener("DOMContentLoaded", function () {
        const finalSection = document.querySelector('.final');
        if (finalSection) {
            const msg = document.createElement("section");
            msg.id = "final";
            msg.setAttribute("class", "final");
            msg.setAttribute("role", "region");
            msg.setAttribute("aria-label", "Submission Confirmation");
            msg.innerHTML = `
              <h2>Your interest is received.</h2>
              <p>We’ll keep you updated on the next batch.</p>
              <a class="to-top" onclick="window.scrollTo({top: 0, behavior: 'smooth'})">
                <span><svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" stroke="#666" stroke-width="2"
                    stroke-linecap="round" stroke-linejoin="round" viewBox="0 0 24 24">
                    <polyline points="18 15 12 9 6 15"></polyline>
                  </svg>
                </span>
                <span class="a-text">To the top</span>
              </a>
            `;
            finalSection.replaceWith(msg);
            msg.scrollIntoView({ behavior: 'smooth' });
        }
    });
}

document.querySelector('.cta').addEventListener('click', function (e) {
    e.preventDefault();
    document.querySelector('.final').scrollIntoView({ behavior: 'smooth' });
});

const form = document.getElementById('waitlist-form');
const submitBtn = form.querySelector('button[type="submit"]');
const inputs = form.querySelectorAll('input[required], select[required]');

function validateEmail(email) {
    const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return re.test(email);
}

function checkFormValidity() {
    let isValid = true;
    inputs.forEach(input => {
        if (input.type === 'email') {
            if (!validateEmail(input.value.trim())) isValid = false;
        } else if (input.value.trim().length < 3) {
            isValid = false;
        }
    });
    submitBtn.disabled = !isValid;
    submitBtn.classList.toggle('disabled', !isValid);
}

inputs.forEach(input => input.addEventListener('input', checkFormValidity));
checkFormValidity();