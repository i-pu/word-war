<template>
  <div id="app">
    <b-navbar>
      <template slot="start">
        <b-navbar-item tag="router-link" :to="{ path: '/' }">
          Top
        </b-navbar-item>
        <b-navbar-item tag="router-link" :to="{ path: '/about' }">
          About
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
    <router-view />
  </div>
</template>

<style>
/*
 * #app {
 *   font-family: "Avenir", Helvetica, Arial, sans-serif;
 *   -webkit-font-smoothing: antialiased;
 *   -moz-osx-font-smoothing: grayscale;
 *   text-align: center;
 *   color: #2c3e50;
 * }
 * 
 * #nav {
 *   padding: 30px;
 * }
 * 
 * #nav a {
 *   font-weight: bold;
 *   color: #2c3e50;
 * }
 * 
 * #nav a.router-link-exact-active {
 *   color: #42b983;
 * }
 */
</style>
