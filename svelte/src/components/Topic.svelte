<script>
  import Card, {Content, PrimaryAction, Media, MediaContent, Actions, ActionButtons, ActionIcons} from '@smui/card';
  import Button, {Label, Icon} from '@smui/button';
  import IconButton from '@smui/icon-button';

  import { createEventDispatcher, onMount } from 'svelte';

  export let topic = '', banner = '', question = '', prediction = '', answer = '', context = '', state = 'published', closing

  const dispatch = createEventDispatcher();
  const localTimeFormat = { 
    weekday: 'long', 
    year: 'numeric', 
    month: 'long', 
    day: 'numeric',
    hour: 'numeric',
    minute: 'numeric'
  }

  let info = ''
  let action = 'Tebak!'
  let disabled = false

  onMount(async () => {
    let now = new Date().getTime()
    let close = new Date(closing)
    let closed = now >= close.getTime()

    if (state !== 'answered' && closed) {
      state = "closed"
    }


    if (prediction) {
      action = "Subah ditebak"
      disabled = true
    }

    switch(state) {
      case "closed":
        info = "Tebak-tebakan ditutup"
        action = "Sudah Tutup"
        disabled = true
        break
      case "answered":
        info = "Tebakan sudah terjawab"
        action = "Terjawab"
        disabled = true
        break
      default:
        info = "Tutup hari " + close.toLocaleDateString("id-ID", localTimeFormat)
        break
    }
  })

  function tebak() {
    if (disabled) return;

    dispatch('tebak', {
      topic: JSON.parse(JSON.stringify(topic))
    })
  }
</script>

<style>
  h2, h3, h4, h5, h6 {
    margin: 0.25rem;
  }

  .title {
    font-weight: bold;
  }

  .subtitle {
    font-weight: normal;
    color: dimgrey;
  }

  .col {
    text-align: right;
  }

  h2.ans:empty::before {
    content: '-';
  }

  .gr-2 {
    display: grid;
    grid-template-columns: 1fr 1fr;
  }
</style>

<Card>
  <Content>
    <h3 class="title">{question}</h3>
    <h5 class="subtitle">{context}</h5>
    <h5 class="subtitle">{info}</h5>
    <div class="gr-2">
    
      <span class="col">
        <h5 class="subtitle">Tebakan</h5>
        <h2 class="ans">{prediction}</h2>
      </span>
      <span class="col">
        <h5 class="subtitle">Jawaban</h5>
        <h2 class="ans">{answer}</h2>
      </span>
    </div>
  </Content>
  <Actions>
    <!-- TODO: Pending implementation
    <ActionButtons>
      <IconButton class="material-icons" on:click={() => alert("do something!")} title="Share">share</IconButton>
      <IconButton on:click={() => alert("do something!")} toggle aria-label="Add to favorites" title="Add to favorites">
        <Icon class="material-icons" on>favorite</Icon>
        <Icon class="material-icons">favorite_border</Icon>
      </IconButton>
    </ActionButtons> -->

    <ActionIcons>
      <Button disabled={disabled} on:click={tebak}>
        <Label>{action}</Label>
      </Button>
    </ActionIcons>

  </Actions>
</Card>