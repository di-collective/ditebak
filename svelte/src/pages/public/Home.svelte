<script>
  // router
  // import { Route } from 'svelte-router-spa'

  // components import
  import Dialog, { Title, Content, Actions, InitialFocus } from '@smui/dialog'
  import Button, { Label } from '@smui/button'
  import Chip, {Set, Icon, Checkmark, Text} from '@smui/chips'
  import Textfield, { Input } from '@smui/textfield'
  import HelperText from '@smui/textfield/helper-text/index'

  import Topic from 'components/Topic.svelte'
  import Gap from 'components/Gap.svelte'
  
  // library imports
  import { onMount } from 'svelte'
  import { writable } from 'svelte/store'
  import { gateway } from 'gateway-conf/gateway.js'

  export let currentRoute, params

  let popup = {
    window: {},
    title: '',
    message: '',
    show: (title, message) => {
      popup.title = title
      popup.message = message

      popup.window.open()
    }
  }
  let dialog = {
    popup: {},
    topic: {},
    bet: {
      topic: '',
      prediction: '',
      stake: 1
    }
  }
  let pagination = {
    next: () => {
      pagination.bind(pagination.page+1)
    },
    enabled: true, loading: false,
    page: 1, total: 0,
    text: '',
    paging: (page, {total_data, total_page}) => {
      pagination.page = page
      pagination.total = total_page
      pagination.enabled = page < total_page
      pagination.text = pagination.enabled ? 'Tampilkan Lebih Banyak' : 'Sudah Habis'
    },
    bind: async (page = 1) => {
      try {
        pagination.loading = true
        const response = await gateway.topics(page)
        const body = response.data
        const more = body.data.map(e => {
          return e
        })

        pagination.paging(page, body.paging)
        topics = [...topics, ...more]
      } catch(e) {
        console.error("Failed to fetch topics", e)
      } finally {
        pagination.loading = false
      }
    }
  }
  let topics = []

  function tebak(event) {
    dialog.topic = event.detail.topic
    dialog.question = event.detail.topic.question
    dialog.popup.open()
  }

  async function submit() {
    try {
      await gateway.bet({
        topic: dialog.topic.id,
        prediction: dialog.bet.prediction.toString(),
        stake: dialog.bet.stake
      })
      popup.show("Sukses nebak!", "Tebakanmu sukses tercatat")
    }
    catch(err) {
      console.error("Failed to place a bet: ", JSON.stringify(err))
      if (err.status === 409) {
        popup.show("Loh? udah ada!", "Anda sudah pernah nebak di topik ini")
        return
      }

      if (err.status === 400) {
        popup.show("Duh, gagal nebak!", "Pastikan tebakan terisi, dan reputasi yang dipertaruhkan antara 1 sampai 10")
      }

      console.error("Failed to place a bet: ", e)
    }
  }

  onMount(async () => {
    pagination.bind()
  })
</script>

<style>
  .chip-right {
    flex-direction: row-reverse;
  }
</style>

<div>
  {#each topics as topic, idx}
  <Topic 
    topic={topic}
    banner={topic.banner}
    question={topic.question}
    prediction={gateway.myBet(topic.id)}
    answer={topic.answer}
    context={topic.context}
    state={topic.state}
    closing={topic.closing_at}
    on:tebak={tebak}>
  </Topic>
  <Gap />
  {/each}

  <Button
    variant="unelevated"
    on:click={pagination.next}
    disabled="{pagination.loading || pagination.page >= pagination.total}">
    {pagination.text}
  </Button>

  <Gap />

  <Dialog bind:this={popup.window} aria-labelledby="popup-title" aria-describedby="popup-message">
    <Title id="popup-title">{popup.title}</Title>
    <Content id="popup-message">
      <Text>{popup.message}</Text>
    </Content>
  </Dialog>

  <Dialog bind:this={dialog.popup} aria-labelledby="tebak-title" aria-describedby="tebak-content">
    <Title id="tebak-title">Tebak</Title>
    <Content id="simple-content">
      <Text>{dialog.topic.question}</Text>

      <Textfield type="number" use={[InitialFocus]} style="width: 100%" dense bind:value={dialog.bet.prediction} label="Tebakanmu" input$aria-controls="helper-text" input$aria-describedby="helper-text-dense" />
      <HelperText id="helper-text">Tulis tebakanmu dan 'Yok Lah!'</HelperText>

      <Textfield type="number" input$min="1" input$max="10" style="width: 100%" dense bind:value={dialog.bet.stake} label="Reputasi" input$aria-controls="helper-text" input$aria-describedby="helper-text-dense" />
      <HelperText id="helper-text">Berapa reputasi yang mau anda pertaruhkan (1 - 10)?</HelperText>
    </Content>
    <Actions>
      <Button color="secondary">
        <Label>Gak Jadi...</Label>
      </Button>
      <Button default on:click={submit}>
        <Label>Yok Lah!</Label>
      </Button>
    </Actions>
  </Dialog>

</div>