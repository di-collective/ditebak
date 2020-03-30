import PublicLayout from './pages/public/Layout.svelte'
import Home from './pages/public/Home.svelte'
import Login from './pages/public/Login.svelte'

const routes = [{
    name: '/',
    layout: PublicLayout,
    component: Home
}, {
    name: 'login',
    layout: PublicLayout,
    component: Login
}]

export { routes }