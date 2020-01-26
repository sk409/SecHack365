import colors from "vuetify/es5/util/colors";

export default {
  mode: "spa",
  /*
   ** Headers of the page
   */
  head: {
    titleTemplate: "%s - " + process.env.npm_package_name,
    title: process.env.npm_package_name || "",
    meta: [{
        charset: "utf-8"
      },
      {
        name: "viewport",
        content: "width=device-width, initial-scale=1"
      },
      {
        hid: "description",
        name: "description",
        content: process.env.npm_package_description || ""
      }
    ],
    link: [{
      rel: "icon",
      type: "image/x-icon",
      href: "/favicon.ico"
    }],
    script: [{
        src: "https://cdnjs.cloudflare.com/ajax/libs/ace/1.2.0/ace.js"
      },
      {
        src: "https://cdnjs.cloudflare.com/ajax/libs/ace/1.2.0/ext-language_tools.js"
      }
    ]
  },
  /*
   ** Customize the progress-bar color
   */
  loading: {
    color: "#fff"
  },
  /*
   ** Global CSS
   */
  css: [{
    src: "@/assets/sass/global.scss",
    lang: "scss"
  }],
  /*
   ** Plugins to load before mounting the App
   */
  plugins: ["@/plugins/filters.js", "@/plugins/mavon_editor.js", "@/plugins/routes.js", "@/plugins/user.js", "@/plugins/utils.js"],
  /*
   ** Nuxt.js dev-modules
   */
  buildModules: ["@nuxtjs/vuetify"],
  /*
   ** Nuxt.js modules
   */
  modules: [],
  /*
   ** vuetify module configuration
   ** https://github.com/nuxt-community/vuetify-module
   */
  vuetify: {
    customVariables: ["~/assets/variables.scss"],
    theme: {
      themes: {
        light: {
          primary: "#03a9f4",
          accent: "#ff7043",
          secondary: "#00bcd4",
          info: colors.teal.lighten1,
          warning: colors.amber.base,
          error: colors.deepOrange.accent4,
          success: colors.green.accent3
        }
      }
    }
  },
  /*
   ** Build configuration
   */
  build: {
    /*
     ** You can extend webpack config here
     */
    extend(config, ctx) {}
  },
  server: {
    port: 3000,
    host: "0.0.0.0"
  },
  env: {
    serverHost: "http://localhost",
    serverPort: "3300",
    serverOrigin: "http://localhost:3300"
  }
};
