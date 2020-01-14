<template>
  <div ref="ide" class="h-100">
    <v-container fluid class="h-100 pa-0">
      <v-row class="h-100">
        <v-col cols="2" class="h-100 pa-0 pl-3">
          <FileTree
            :items="filetree"
            :load-children="fetchChildren"
            :style="filetreeStyle"
            class="file-tree"
            @click-item="clickFileTreeItem"
          ></FileTree>
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
import _ from "lodash";
import ajax from "@/assets/js/ajax.js";
import FileTree from "@/components/FileTree.vue";
import { Url, urlFiles, urlFolders, urlLessons } from "@/assets/js/url.js";

const extension = str => {
  const components = str.split(".");
  const extension =
    components.length == 0 ? "" : components[components.length - 1];
  return extension;
};

const aceMode = path => {
  const ext = extension(path);
  const modes = {
    js: "javascript",
    php: "php"
  };
  const base = "ace/mode/";
  if (!modes[ext]) {
    return base + "text";
  }
  return base + modes[ext];
};

const fileIcon = name => {
  const ext = extension(name);
  if (
    ext == "bash_history" ||
    ext == "bash_logout" ||
    ext == "bash_profile" ||
    ext == "bashrc"
  ) {
    return "mdi-bash";
  } else if (ext == "php") {
    return "mdi-language-php";
  }
  return "mdi-note-text-outline";
};

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
        c.icon = fileIcon(child.Name);
      } else {
        c.file = false;
        c.open = false;
        c.icon = "mdi-folder";
        c.children = [];
      }
      return c;
    }
  );
};

let editor = null;
export default {
  layout: "lesson",
  components: {
    FileTree
  },
  data() {
    return {
      delayedUpdate: _.debounce(this.updateFileText, 1000),
      file: null,
      filetree: [],
      lesson: null,
      ideHeight: 0
    };
  },
  computed: {
    consoleURL() {
      return this.lesson
        ? process.env.serverHost + ":" + this.lesson.HostConsolePort
        : "";
    },
    filetreeStyle() {
      return {
        height: this.ideHeight + "px"
      };
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
    this.ideHeight = this.$refs.ide.clientHeight;
    this.setupAce();
  },
  methods: {
    clickFileTreeItem(item) {
      if (item.file) {
        const url = new Url(urlFiles);
        const data = {
          lessonID: this.lesson.ID,
          path: item.path
        };
        ajax.get(url.base, data).then(response => {
          item.text = response.data.Text;
          this.file = item;
          editor.setValue(item.text);
          editor.getSession().setMode(aceMode(item.path));
        });
      } else {
        item.open = !item.open;
        item.icon = item.open ? "mdi-folder-open" : "mdi-folder";
        if (!item.open) {
          item.children = [];
        }
      }
    },
    fetchChildren(item) {
      if (item.file) {
        return;
      }
      const url = new Url(urlFolders);
      const data = {
        lessonID: this.$route.params.id,
        path: item.path
      };
      ajax.get(url.base, data).then(response => {
        item.children = makeChildrenFromResponse(response);
      });
    },
    setupAce() {
      ace.edit("editor");
      editor = ace.edit("editor");
      editor.$blockScrolling = Infinity;
      editor.setOptions({
        enableBasicAutocompletion: true,
        enableSnippets: false,
        enableLiveAutocompletion: true
      });
      editor.setTheme("ace/theme/xcode");
      editor.setFontSize(20);
      editor.session.on("change", delta => {
        if (delta.action === "insert" && this.file.text !== editor.getValue()) {
          this.file.text = editor.getValue();
          this.delayedUpdate();
        }
      });
    },
    updateFileText() {
      const url = new Url(urlFiles);
      const data = {
        lessonID: this.lesson.ID,
        path: this.file.path,
        text: this.file.text
      };
      ajax.put(url.base, data);
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
}
.file-tree {
  overflow: scroll;
  white-space: nowrap;
}
</style>
