<template>
  <div class="h-100">
    <v-container fluid class="h-100 pa-0">
      <v-row class="h-100">
        <v-col cols="2" class="h-100 pa-0">
          <v-treeview :items="filetree" :load-children="fetchChildren"></v-treeview>
        </v-col>
        <v-col cols="10" class="h-100 pa-0">
          <div id="editor" class="h-100"></div>
          <iframe :src="consoleURL" class="console"></iframe>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import { Url, urlFolders, urlLessons } from "@/assets/js/url.js";

const makeChildrenFromResponse = response => {
  const children = [];
  if (!response.data.ChildFolders) {
    response.data.ChildFolders = [];
  }
  if (!response.data.ChildFiles) {
    response.data.ChildFiles = [];
  }
  return response.data.ChildFolders.concat(response.data.ChildFiles).map(
    child => {
      const c = {
        name: child.Name,
        path: child.Path
      };
      if (child.hasOwnProperty("Text")) {
        c.file = true;
      } else {
        c.file = false;
        c.children = [];
      }
      return c;
    }
  );
};
export default {
  layout: "lesson",
  data() {
    return {
      filetree: [],
      lesson: null
    };
  },
  computed: {
    consoleURL() {
      return this.lesson
        ? process.env.serverHost + ":" + this.lesson.HostConsolePort
        : "";
    }
  },
  created() {
    const lessonID = this.$route.params.id;
    const urlL = new Url(urlLessons);
    const lessonData = {
      id: lessonID
    };
    ajax.get(urlL.base, lessonData).then(response => {
      this.lesson = response.data[0];
    });
    const urlF = new Url(urlFolders);
    const folderData = {
      lessonID,
      path: "/"
    };
    ajax.get(urlF.base, folderData).then(response => {
      this.filetree = makeChildrenFromResponse(response);
    });
  },
  mounted() {
    ace.edit("editor");
    const editor = ace.edit("editor");
    editor.$blockScrolling = Infinity;
    editor.setOptions({
      enableBasicAutocompletion: true,
      enableSnippets: false,
      enableLiveAutocompletion: true
    });
    editor.setTheme("ace/theme/xcode");
    editor.getSession().setMode("ace/mode/javascript");
    editor.setFontSize(20);
    editor.setValue("console.log()");
  },
  methods: {
    fetchChildren(item) {
      const url = new Url(urlFolders);
      const data = {
        lessonID: this.$route.params.id,
        path: item.path
      };
      ajax.get(url.base, data).then(response => {
        item.children = makeChildrenFromResponse(response);
      });
    }
  }
};
</script>

<style>
#editor {
  height: 60%;
}
.console {
  width: 100%;
  height: 40%;
  /* background: black; */
}
</style>
