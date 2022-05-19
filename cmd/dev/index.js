const app = Vue.createApp({
    delimiters: ['${', '}']
})

const List = {
    template: '<basic-list></basic-list>'
}

const Item = {
    template: '<basic-item></basic-item>'
}

const AddBasic = {
    template: '<basic-make></basic-make>'
}

const NotFoundComponent = {
    template: '<div><h2>Page not found</h2></div>'
}

//-{{template "cmp-js"}}

const router = new VueRouter.createRouter({
    history: VueRouter.createWebHistory(),
    routes: [
        { path: '/', component: List},
        { path: '/basics', component: List},
        { path: '/basics/', component: List},
        { path: '/basics/:uc', component: Item},
        { path: '/basics/add', component: AddBasic},
        { path: '/:pathMatch(.*)', component: NotFoundComponent }
    ]
})
app.use(router)

const store = Vuex.createStore({
    state () { return {
        //-{{template "cmp-store-state"}}
    }},
    mutations: {
        //-{{template "cmp-store-mutations"}}
    },
    actions: {
        //-{{template "cmp-store-actions" .}}
    }
})
app.use(store)

app.mount('#app')

store.dispatch('readBasics');