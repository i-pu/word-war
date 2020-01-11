<template>
  <div id="app">
    <b-navbar>
      <template slot="start">
        <b-navbar-item tag="router-link" :to="{ path: '/' }">
          Top
        </b-navbar-item>
        <b-navbar-item tag="router-link" :to="{ path: '/home' }">
          Home
        </b-navbar-item>
        <b-navbar-item tag="router-link" :to="{ path: '/game' }">
          Game
        </b-navbar-item>
        <b-navbar-item tag="router-link" :to="{ path: '/result' }">
          Result
        </b-navbar-item>
      </template>
    </b-navbar>
  </div>
</template>
