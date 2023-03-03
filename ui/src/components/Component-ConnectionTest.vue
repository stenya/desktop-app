<template>
  <div class="flexColumn">
    <div class="flexColumn" style="margin: 20px">
      <div class="large_text" style="text-align: center; margin-top: 130px">
        {{ title }}
      </div>

      <div style="text-align: center; margin-top: 20px">
        <div class="flexRow">
          <div class="statusName">Protocol:</div>
          <div class="statusValue">WireGuard</div>
        </div>
        <div class="flexRow">
          <div class="statusName">Server:</div>
          <div class="statusValue">{{ statusServer }}</div>
        </div>
        <div class="flexRow">
          <div class="statusName">Host:</div>
          <div class="statusValue">{{ statusHost }}</div>
        </div>
        <div class="flexRow">
          <div class="statusName">Port:</div>
          <div class="statusValue">{{ statusPort }}</div>
        </div>
      </div>

      <!--
      <div style="background: lightgreen; margin-top: 20px; height: 100%"></div>
      -->
      <div style="margin-top: 50px; margin-bottom: 20px">
        <div>
          <div>
            <button
              v-if="isFinishedSuccess"
              class="master"
              v-on:click="Apply"
              style="margin-bottom: 10px"
            >
              Apply & Connect
            </button>
            <button class="slave" v-on:click="Cancel">Cancel</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { VpnTypeEnum } from "@/store/types";
const sender = window.ipcSender;

async function connect() {
  try {
    await sender.Connect();
  } catch (e) {
    console.error(e);
    sender.showMessageBoxSync({
      type: "error",
      buttons: ["OK"],
      message: `Failed to connect: ` + e,
    });
  }
}

export default {
  components: {},
  data: function () {
    return {};
  },
  mounted() {},
  computed: {
    isRunning: function () {
      return this.$store.state.vpnState.connectionTestRunning;
    },
    isFinishedSuccess: function () {
      let found = this.$store.state.vpnState.connectionTestResult;
      return (
        !this.isRunning &&
        found != null &&
        found.GoodConnection != null &&
        !found.Error
      );
    },
    title: function () {
      if (this.isFinishedSuccess)
        return "Success! Workable connection parameters were found.";
      if (this.isRunning) return "Testing servers connectivity";

      return "No connectivity found";
    },
    ctStatus: function () {
      return this.$store.state.vpnState.connectionTestRunningState;
    },
    statusServer: function () {
      if (!this.ctStatus || !this.ctStatus.Server) return "...";
      return `${this.ctStatus.Server.city} (${this.ctStatus.Server.country_code})`;
    },
    statusHost: function () {
      if (!this.ctStatus || !this.ctStatus.Host || !this.ctStatus.Port)
        return "...";

      return `${this.ctStatus.Host.hostname} `;
    },
    statusPort: function () {
      if (!this.ctStatus || !this.ctStatus.Port) return "...";

      return ` ${this.ctStatus.Port.port}:${this.ctStatus.Port.type}`;
    },
  },
  watch: {},
  methods: {
    Cancel() {
      sender.ConnectionTest_Stop();
      this.$store.dispatch(`uiState/isConnectionTestView`, false);
    },

    Apply() {
      let found =
        this.$store.state.vpnState.connectionTestResult.GoodConnection;

      this.$store.dispatch(`settings/vpnType`, VpnTypeEnum.WireGuard);

      let svr = this.$store.state.vpnState.serversHashed[found.Gateway];
      if (!svr) {
        console.error("Server not found:", found);

        console.log(this.$store.state.vpnState.serversHashed);
        return;
      }
      this.$store.dispatch(`settings/serverEntry`, svr);

      this.$store.dispatch("settings/showHosts", true);
      this.$store.dispatch(`settings/serverEntryHostId`, found.HostName);

      this.$store.dispatch(`settings/setPort`, {
        port: found.Port,
        type: found.PortType,
      });
      this.$store.dispatch(`uiState/isConnectionTestView`, false);

      connect();
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
@import "@/components/scss/constants";

.large_text {
  margin: 12px;
  font-weight: 600;
  font-size: 18px;
  line-height: 120%;
}

.small_text {
  margin: 12px;
  margin-top: 0px;

  font-size: 13px;
  line-height: 17px;
  letter-spacing: -0.208px;

  color: #98a5b3;
}
.statusName {
  opacity: 0.7;
  text-align: left;
  min-width: 60px;
  font-size: 13px;
}
.statusValue {
  font-size: 13px;
}
</style>
