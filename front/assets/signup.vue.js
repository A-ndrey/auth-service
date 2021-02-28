const Signup = {
    template: `
        <div class="form">
        <h2>Signup</h2>
        <h3>{{ service }}</h3>
        <div class="input">
            <div class="inputBox">
                <label>Email</label>
                <input type="email" v-model="email" placeholder="example@email.com">
                <p class="message" :class="emailStatus">{{ emailMessage }}</p>
            </div>
            <div class="inputBox">
                <label>Password</label>
                <input type="password" v-model="password" placeholder="••••••••">
                <p class="message" :class="passwordStatus">{{ passwordMessage }}</p>
            </div>
            <div class="inputBox">
                <label>Confirm Password</label>
                <input type="password" v-model="confirmPassword" placeholder="••••••••">
                <p class="message" :class="confirmPasswordStatus">{{ confirmPasswordMessage }}</p>
            </div>
            <div class="inputBox">
                <input type="submit" value="Sign Up" @click="signup">
            </div>
        </div>
        <p class="switch-form">Already have an account? <span class="link" @click="signin">Sign In</span></p>
    </div>
    `,
    data: () => ({
        email: "",
        password: "",
        confirmPassword: "",
        emailMessage: "",
        emailStatus: "",
        passwordMessage: "",
        passwordStatus: "",
        confirmPasswordMessage: "",
        confirmPasswordStatus: "",
    }),
    methods: {
        clearValidation() {
            this.emailMessage = ""
            this.emailStatus = ""
            this.passwordMessage = ""
            this.passwordStatus = ""
            this.confirmPasswordMessage = ""
            this.confirmPasswordStatus = ""
        },
        signin() {
            this.$root.content = "Signin"
        },
        async signup() {
            this.clearValidation()

            if (this.password !== this.confirmPassword) {
                this.confirmPasswordMessage = 'Passwords should not be different'
                this.confirmPasswordStatus = 'status-error'
                return false
            }

            fetch(`${this.apiURL}/signup?service=${this.service}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    email: this.email,
                    password: this.password
                })
            })
                .then(response => response.json())
                .then(data => {
                    if (data.error) {
                        if (data.error_type === 'email') {
                            this.emailMessage = data.error
                            this.emailStatus = 'status-error'
                        } else if (data.error_type === 'password') {
                            this.passwordMessage = data.error
                            this.passwordStatus = 'status-error'
                        } else if (data.error_type === 'common') {
                            this.confirmPasswordMessage = data.error
                            this.confirmPasswordStatus = 'status-error'
                        }
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
        }
    },
    computed: {
        service() {
            return this.$root.service
        },
        apiURL() {
            return this.$root.apiURL
        }
    }
}