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
    reputation: profile.reputation,
    bets: []
  })
  await syncMyBets()
}

async function doLogout() {
  await rest.get(config.logout)
  person.toStorage({
    displayName: '',
    email: '',
    photo: '',
    bets: []
  })
}

async function placeBet(command) {
  const response = await rest.post(config.bets, {
    data: command
  })
  const body = response.data
  const bet = body.data
  person.addBets(bet)
}

async function syncMyBets() {
  const response = await rest.get(config.bets)
  const body = response.data
  const bets = body.data
  person.addBets(...bets)
}

async function syncProfile() {
  const response = await rest.get(config.profile)
  const body = response.data
  const profile = body.data
  
  person.toStorage({
    displayName: profile.display_name,
    email: profile.email,
    photo: profile.photo,
    reputation: profile.reputation,
    bets: []
  })
  await syncMyBets()
}

const gateway = {
  login: doLogin,
  logout: doLogout,
  topics: findAllTopics,
  bet: placeBet,

  myBet: person.myBet,
  sync: syncProfile
}

export { gateway }