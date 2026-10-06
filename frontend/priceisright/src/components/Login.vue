<template>
  <div>
    <StarField />
    <h1 class="text-center h4 text-uppercase">
      <img
        class="titleBanner mt-2"
        src="@/assets/images/logo.png"
        alt="Price is Right"
        title="Price is Right"
        height="1564"
        width="1414"
      />
    </h1>
    <form class="ml-auto mr-auto mb-4" @submit.prevent="connectGame">
      <fieldset class="rounded-lg p-2">
        <legend class="text-center">Join the Fun!</legend>
        <div class="form-row">
          <div class="form-group col">
            <input
              v-model="Name"
              class="form-control"
              placeholder="Your Name"
              aria-label="Your Name"
              required
            />
          </div>
        </div>
        <div class="form-row">
          <div class="form-group col">
            <input
              v-model="BadgeNumber"
              class="form-control"
              placeholder="Badge Number"
              aria-label="Badge Number"
              aria-describedby="badgeNumberHelp"
              pattern="[0-9]+"
              required
            />
            <small id="badgeNumberHelp" class="form-text text-white">
              <!-- Include badge letter. E.g. 123A. -->
              Please do not join unless you're sat in the main stage audience.
            </small>
          </div>
        </div>
        <div class="form-row">
          <div class="form-group col text-center mb-0">
            <button class="btn btn-primary">Connect to Game</button>
          </div>
        </div>
        <div v-if="Errors.Connect" class="mt-2">
          <span class="form-text text-danger text-center text-capitalize">
            {{ Errors.Connect }}
          </span>
        </div>
      </fieldset>
    </form>
  </div>
</template>
<script>
import Axios from "axios";
import StarField from "./backgrounds/StarField.vue";

export default {
  name: "LoginView",
  components: {
    StarField,
  },
  data() {
    return {
      Errors: {
        Connect: "",
      },
      Name: "",
      BadgeNumber: "",
      Password: "",
    };
  },
  mounted() {
    const urlParams = new URLSearchParams(window.location.search);
    if (urlParams.has("action")) {
      if (!urlParams.has("name")) {
        return;
      }
      if (!urlParams.has("badge")) {
        return;
      }
      this.Name = urlParams.get("name");
      this.BadgeNumber = urlParams.get("badge");
      switch (urlParams.get("action")) {
        case "join":
          this.connectGame();
          break;
        case "control":
          if (!urlParams.has("password")) {
            return;
          }
          this.Password = urlParams.get("password");
          this.connectGame();
          break;
        default:
      }
    }
  },
  methods: {
    connectGame() {
      this.Errors.Connect = "";
      const data = {
        BadgeNumber: this.BadgeNumber,
        Name: this.Name,
        Password: this.Password,
      };
      Axios({
        method: "post",
        url: `${this.apiUri}/game/join`,
        data,
      })
        .then((r) => {
          let game;
          if (r.data.Response.Game !== undefined) {
            game = r.data.Response.Game;
          }
          const player = r.data.Response.Player;
          const token = r.data.Response.Tokens;

          token.Endpoint = `${this.apiUri}/token/refresh`;
          token.Service = "PriceIsRight";
          token.Vendor = this.RoomCode;

          this.emitter.emit("newToken", token);
          this.$emit("joinedGame", player, game);
        })
        .catch((e) => {
          if (typeof e.response === "object") {
            this.Errors.Connect = e.response.data.ErrorMessage;
          }
        });
    },
  },
};
</script>
<style scoped>
form {
  max-width: 350px;
}
legend {
  position: absolute;
  left: -10000px;
}
.titleBanner {
  height: auto;
  width: 250px;
  max-width: 90%;
}
</style>
