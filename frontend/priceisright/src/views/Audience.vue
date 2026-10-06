<template>
  <div>
    <StarField />
    <div class="main-display min-vh-100">
      <h1 class="bg-dark position-relative text-center">
        <a class="position-absolute text-white pl-2 pr-2"
           href="#"
           title="Disconnect"
           @click="$emit('quitGame')"
        >
          <FontAwesomeIcon :icon="['fas', 'xmark']" />
        </a>
        Shadow's Price Is Right
      </h1>
      <div class="m-auto">
        <div v-if="gameCode === 0">
          &nbsp;
        </div>
        <div v-else-if="gameCode === 100">
          <DisplayMessage effect=""
                          message="Thank you for watching!"
                          class="d-block p-2"
          />
          <DisplayMessage effect=""
                          message="Shadow Raccoon's Price is Right!"
                          class="d-block p-2"
          />
          <p class="text-center p-2">
            We hope you enjoyed the show.
          </p>
          <p class="text-center p-2">
            <a href="https://twitter.com/shadow_raccoon">Shadow Raccoon</a>
            <br><br>
            <a href="https://twitter.com/FrostyWuffle">Frío Wuffle</a>
            <br><br>
            <a href="https://twitter.com/FreddyPanda">Freddy Panda</a>
            <br><br>
            <a href="https://twitter.com/TheadoreRabbit">Theadore Rabbit</a>
            <br><br>
            <a href="https://twitter.com/Regdeh">Regdeh</a>
            <br><br>
            See you next time!
          </p>
        </div>
        <div v-else-if="gameCode === 1">
          <DisplayMessage effect=""
                          message="You're a bit early!"
                          class="d-block p-2"
          />
          <DisplayMessage effect=""
                          message="Join us for Shadow Raccoon's Price is Right at Scotiacon"
                          class="d-block p-2"
          />
          <DisplayMessage effect=""
                          message="Saturday 3rd Feb 4:00pm"
                          class="d-block p-2"
          />
        </div>
        <div v-else-if="gameCode === 2">
          <DisplayMessage effect=""
                          message="Shadow Raccoon's Price is Right!"
                          class="d-block p-2"
          />
          <DisplayMessage effect=""
                          message="The show will be starting shortly"
                          class="d-block p-2"
          />
          <p class="text-center p-2">
            This is an optional audience participation game
          </p>
          <p class="text-center p-2">
            If you are not in the main stage
            audience or do not wish to be potentially called up on to the main stage,
            please click the x in the top right corner.
          </p>
        </div>
        <NewGuess v-else-if="gameCode === 20" />
        <div v-else-if="gameCode === 30">
          <DisplayMessage effect="ripple"
                          message="Round Closed"
                          class="d-block"
          />
          <p class="text-center p-2">
            You are no longer able to submit guesses for this round
          </p>
          <p class="text-center p-2">
            The closest participants will be called up onto stage at the end of the current game
          </p>
          <p class="text-center p-2">
            Now sit back, enjoy the show and, continue helping
            the contestants by shouting out advice
          </p>
        </div>
        <div v-else-if="gameCode === 40">
          <template v-if="calling.UUID === player.UUID">
            <DisplayMessage effect="ripple"
                            message="Come on Down!"
                            class="d-block"
            />
            <p class="text-center p-2">
              We'd like you to join us on stage
            </p>
          </template>
          <template v-else>
            <DisplayMessage effect="ripple"
                            message="Look up!"
                            class="d-block"
            />
            <p class="text-center p-2">
              It could be your turn to join us on stage
            </p>
          </template>
        </div>
        <div v-else>
          <DisplayMessage effect="ripple"
                          message="Please Wait"
                          class="d-block"
          />
          <p class="text-center">
            During contestants row prompts will appear to allow you to submit your price guesses.
          </p>
          <p class="text-center">
            If no prompt appears, refresh the page.
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import DisplayMessage from '../components/display/DisplayMessage.vue';
import NewGuess from '../components/audience/NewGuess.vue';
import PlayerView from './Player.vue';
import StarField from '../components/backgrounds/StarField.vue';

export default {
  name: 'AudienceView',
  components: {
    FontAwesomeIcon,
    DisplayMessage,
    NewGuess,
    StarField,
  },
  extends: PlayerView,
  props: {
    suggestionsEnabled: {
      type: Boolean,
      required: false,
      default() {
        return false;
      },
    },
  },
};

</script>

<style scoped>
h1 {
  font-family: 'LEMON MILK', sans-serif;
  font-size: 1.5rem !important;
  padding-left: 28px;
  padding-right: 28px;
}
h1 svg {
  margin-top: -4px;
}
h1 a {
  border-left: 4px solid black;
  right: 0 !important;
}
h1 a:hover {
  background-color: lightgray !important;
  color: black !important;
}
h2 {
  color: white;
  font-family: 'LEMON MILK', sans-serif;
  font-size: 1.7rem !important;
  text-shadow:
      -1px -1px 0 #000,
      1px -1px 0 #000,
      -1px 1px 0 #000,
      1px 1px 0 #000;
}
.main-display {
  background-image: url("@/assets/images/logo-faded.png");
  background-size: 50vw;
  background-repeat: no-repeat;
  background-position: center center;
}
</style>
<style>
@media screen and (max-width: 450px) {
  #WebAppBody .main-display h1 {
    font-size: 1.2rem !important;
  }
  h1 svg {
    margin-top: -2px;
  }
}
@media screen and (max-width: 370px) {
  #WebAppBody .main-display h1 {
    font-size: 1rem !important;
  }
  h1 svg {
    margin-top: -2px;
  }
}
</style>
