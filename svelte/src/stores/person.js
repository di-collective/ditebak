import { writable } from 'svelte/store';

const STORAGE_KEY = 'key.user.info'

const state = {
  displayName: '',
  email: '',
  photo: '',
  reputation: 0,
}

const { subscribe, update } = writable(state)

const person = {
  subscribe,
  displayName: (dp) => update(o => {
    o.displayName = dp
    return o
  }),
  email: (e) => update(o => {
    o.email = e
    return o
  }),
  photo: (p) => update(o => {
    o.photo = p
    return o
  }),
  reputation: (r) => update(o => {
    o.reputation = r
    return o
  }),
  toStorage: (info) => update(o => {
    o.displayName = info.displayName
    o.email = info.email
    o.photo = info.photo
    o.reputation = info.reputation

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
    return o
  })
}
person.fromStorage()

export { person }