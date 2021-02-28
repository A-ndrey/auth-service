const Signin = {
    template: `
        <div class="form">
        <h2>Login</h2>
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
                <input type="submit" value="Sign In" @click="signin">
            </div>
        </div>
        <p class="switch-form">Have not an account yet? <span class="link" @click="signup">Sign Up</span></p>
    </div>
    `,
    data: () => ({
        email: "",
        password: "",
        emailMessage: "",
        emailStatus: "",
        passwordMessage: "",
        passwordStatus: ""
    }),
    methods: {
        signup() {
            this.$root.content = "Signup"
        },
        async signin() {
            fetch(`${this.apiURL}/signin?service=${this.service}`, {
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
                        this.password = ""
                        this.passwordMessage = data.error
                        this.passwordStatus = 'status-error'
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