export const service = document.getElementById("service").textContent
export const emailField = document.getElementById("inputEmail")
export const passwordField = document.getElementById("inputPassword")
export const confirmPasswordField = document.getElementById("confirmPassword")
export const submitForm = document.getElementById("submitForm")
export const apiURL = window.location.origin + "/api/v1"

emailField.oninput = () => {
    emailField.classList.remove("is-invalid")
}

if (confirmPasswordField) {
    confirmPasswordField.onchange = () => {
        confirmPasswordField.classList.remove("is-invalid")
    }
}

if (confirmPasswordField) {
    confirmPasswordField.validate = () => {
        if (confirmPasswordField.value !== passwordField.value) {
            confirmPasswordField.classList.add("is-invalid")
            return false
        }
        return true
    }
}
