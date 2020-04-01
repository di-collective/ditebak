import { firebaseConfig } from './settings'

let initialized = false
let Auth, AuthUI, Providers
function init() {
  if (initialized) return {
    AuthUI: AuthUI,
    Providers: Providers
  };
  
  firebase.initializeApp(firebaseConfig)
  initialized = true

  Auth = firebase.auth()
  AuthUI = new firebaseui.auth.AuthUI(Auth),
  Providers = [
    // firebase.auth.EmailAuthProvider.PROVIDER_ID,
    firebase.auth.GoogleAuthProvider.PROVIDER_ID
  ]

  return {
    AuthUI: AuthUI,
    Providers: Providers
  }
}


export { init }
