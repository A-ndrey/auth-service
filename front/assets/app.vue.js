const App = {
    template: `
        <div>
            <transition name="component-fade" mode="out-in">
                <component :is="content()"></component>
            </transition>
        </div>
    `,
    methods: {
        content() {
            return this.$root.content
        }
    }
}

Vue.component("App", App)
Vue.component("Signin", Signin)
Vue.component("Signup", Signup)

const app = new Vue({
    el: "#app",
    data: () => ({
        content: "Signin",
        service: new URL(location.href).searchParams.get("service"),
        apiURL: window.location.origin + "/api/v1",
    }),
    render: h => h(App)
})
