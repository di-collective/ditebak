<script>
  import TopAppBar, {
    Row, 
    Section, 
    Title
  } from '@smui/top-app-bar'
  import IconButton from '@smui/icon-button'
  import Menu, { SelectionGroup, SelectionGroupIcon } from '@smui/menu';
  import List, { Item, Separator, Text, PrimaryText, SecondaryText, Graphic } from '@smui/list';
  import { Anchor } from '@smui/menu-surface'
  import Textfield from '@smui/textfield'
  import Button, { Label } from '@smui/button'

  import { Navigate, navigateTo } from 'svelte-router-spa'
  import { gateway } from '../config/gateway/gateway.js'
  import { person } from '../stores/person.js'

  let menu = {
    profile: {},
    anchor: {}
  }

  function onAccountClicked() {
    gateway.sync()
    menu.profile.setOpen(true)
  }
</script>

<style>
  .pp-img {
    border-radius: 12px;
  }
</style>

<TopAppBar dense>
  <Row>
    <Section>
      <!-- <IconButton class='material-icons'>menu</IconButton> -->
      <Navigate to='/'><Title>Ditebak.com</Title></Navigate>
    </Section>
    <Section align='end' toolbar>
      {#if !$person.photo}
      <Navigate to='login'>
        <IconButton class='material-icons' aria-label='Account'>account_circle</IconButton>
      </Navigate>
      {:else}
      <div use:Anchor bind:this={menu.anchor}>
        <IconButton aria-label='Account' on:click={onAccountClicked}>
          <img class='pp-img' src='{$person.photo}' alt='profile'>
        </IconButton>
      </div>

      <Menu bind:this={menu.profile} anchor={false} bind:anchorElement={menu.anchor} anchorCorner='BOTTOM_LEFT'>
        <List twoLine>
          <Item>
            <Text>
              <PrimaryText>{$person.displayName}</PrimaryText>
              <SecondaryText>{$person.email}</SecondaryText>
            </Text>
          </Item>
          <Item>
            <Text>
              <PrimaryText>{$person.reputation}</PrimaryText>
              <SecondaryText>Reputasi anda saat ini</SecondaryText>
            </Text>
          </Item>
          <Item>
            <Text>
              <PrimaryText>Tebak-tebakan</PrimaryText>
              <SecondaryText>Lihat koleksi tebak-tebakan anda</SecondaryText>
            </Text>
          </Item>
          <Separator />
          <Item on:SMUI:action={() => {
            gateway.logout()
            navigateTo('/')
          }}>
            <Text>
              <PrimaryText>Cabuts Guys...</PrimaryText>
              <SecondaryText>Logout? Salah mulu ya?</SecondaryText>
            </Text>
          </Item>
        </List>
      </Menu>
      {/if}
    </Section>
  </Row>
</TopAppBar>