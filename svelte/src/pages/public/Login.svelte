<script>
  import { onMount } from 'svelte'
  import { navigateTo } from 'svelte-router-spa'

  import { AuthUI, Providers } from 'firebase-conf/firebase.js'
  import { gateway } from 'gateway-conf/gateway.js'

  const conf = {
    callbacks: {
      signInSuccessWithAuthResult: function(authResult, redirectUrl) {
        const fa = JSON.parse(JSON.stringify(authResult))
        gateway.login(fa).then(() => {
          navigateTo('/')
        })

        // User successfully signed in.
        // Return type determines whether we continue the redirect automatically
        // or whether we leave that to developer to handle.
        // NOTES: THIS DOESN'T WORK IF WE USE ASYNC FUNCTION / RETURNING PROMISE
        return false
      },
      uiShown: function() {
        // The widget is rendered.
        // Hide the loader.
        document.getElementById('loader').style.display = 'none';
      }
    },
    // Will use popup for IDP Providers sign-in flow instead of the default, redirect.
    signInFlow: 'popup',
    signInOptions: Providers,
    signInSuccessUrl: '<my-redirect-asshole!>',
    // Terms of service url.
    tosUrl: '<your-tos-url>',
    privacyPolicyUrl: '<your-privacy-policy-url>'
  }

  onMount(() => {
    AuthUI.start('#firebaseui-auth-container', conf)
  })
</script>

<div id='firebaseui-auth-container' />
<div id='loader'>Loading...</div>
