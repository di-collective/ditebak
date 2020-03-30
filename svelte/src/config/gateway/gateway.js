import axios from 'axios'
import { navigateTo } from 'svelte-router-spa'

import { config } from './settings'
import { person } from '../../stores/person.js'

const rest = axios.create({
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  },
  withCredentials: true
})

// TODO: intercept 401
rest.interceptors.response.use(response => {
  return response
}, e => {
  if (e.response.status === 401) {
    alert("Login dulu ya sebelum nebak...")
    navigateTo('login')
  }

  return Promise.reject(e.response);
})

async function findAllTopics(page = 1) {
  console.log("PAGE=",page)
  return rest.get(config.topics, {
    params: {
      page: page
    }
  })
}

async function doLogin(fa) {
  const response = await rest.post(config.login, {
    data: fa
  })
  const body = response.data
  const profile = body.data

  person.toStorage({
    displayName: profile.display_name,
    email: profile.email,
    photo: profile.photo,
    reputation: profile.reputation
  })
}

async function doLogout() {
  await rest.get(config.logout)
  person.toStorage({
    displayName: '',
    email: '',
    photo: ''
  })
}

async function placeBet(bet) {
  return await rest.post(config.bets, {
    data: bet
  })
}

async function myBets() {
  const response = await rest.get(config.bets)
  const body = response.data
  return body
}

async function syncProfile() {
  const response = await rest.get(config.profile)
  const body = response.data
  const profile = body.data
  
  person.toStorage({
    displayName: profile.display_name,
    email: profile.email,
    photo: profile.photo,
    reputation: profile.reputation
  })
}

const gateway = {
  login: doLogin,
  logout: doLogout,
  topics: findAllTopics,
  bet: placeBet,

  sync: syncProfile
}

export { gateway }