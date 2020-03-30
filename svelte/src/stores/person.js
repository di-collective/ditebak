import { writable } from 'svelte/store';

const STORAGE_KEY = 'key.user.info'

const state = {
  displayName: '',
  email: '',
  photo: '',
  reputation: 0,
  bets: []
}

const { subscribe, update } = writable(state)

const person = {
  subscribe,
  myBet: (topic) => {
    const bet = state.bets.find((b) => {
      return b.topic_id === topic
    })
    if (!bet) return '';

    return bet.prediction
  },
  addBets: (...b) => update(o => {
    o.bets.push(...b)

    localStorage.setItem(STORAGE_KEY, JSON.stringify(o))
    return o
  }),
  toStorage: (info) => update(o => {
    o.displayName = info.displayName
    o.email = info.email
    o.photo = info.photo
    o.reputation = info.reputation
    o.bets = info.bets

    localStorage.setItem(STORAGE_KEY, JSON.stringify(o))
    return o
  }),
  fromStorage: () => update(o => {
    const stored = localStorage.getItem(STORAGE_KEY)
    if (!stored) return o;

    const info = JSON.parse(stored)
    o.displayName = info.displayName
    o.email = info.email
    o.photo = info.photo
    o.reputation = info.reputation
    o.bets = info.bets
    return o
  })
}
person.fromStorage()

export { person }