import {
    service,
    emailField,
    passwordField,
    confirmPasswordField,
    submitForm,
    apiURL
} from './common.js'

submitForm.onsubmit = () => {
    if (submitForm.checkValidity() === false) {
        submitForm.classList.add('was-validated')
        return false
    }

    if (confirmPasswordField.validate() === false) {
        return false
    }

    fetch(`${apiURL}/signup?service=${service}`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            email: emailField.value,
            password: passwordField.value
        })
    })
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                emailField.classList.add('is-invalid')
            } else {
                const params = new URLSearchParams(location.search)
                const tokens = new URLSearchParams(data)
                const redirectURL = params.get('redirect_url')
                if (window.opener && redirectURL) {
                    window.opener.location.replace(redirectURL + '#' + tokens)
                }
                window.close()
            }
        })
    return false
}
