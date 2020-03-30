import { firebaseConfig } from './settings'

firebase.initializeApp(firebaseConfig)

const Auth = firebase.auth()
const Providers = [
    // firebase.auth.EmailAuthProvider.PROVIDER_ID,
    firebase.auth.GoogleAuthProvider.PROVIDER_ID
]

const AuthUI = new firebaseui.auth.AuthUI(Auth)

export { Auth, AuthUI, Providers }
