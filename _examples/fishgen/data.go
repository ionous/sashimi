package fishgen

var Data = `{
 "actions": {
  "acquire-it": {
   "id": "acquire-it",
   "name": "acquire it",
   "event": "acquiring-it",
   "nouns": [
    "actors",
    "props"
   ],
   "actions": [
    "a38fd7888ecbfb8e121d362dbfed"
   ]
  },
  "attack-it": {
   "id": "attack-it",
   "name": "attack it",
   "event": "attacking-it",
   "nouns": [
    "actors",
    "objects"
   ],
   "actions": [
    "e1067e858bff37ff95ed58744c3655e0"
   ]
  },
  "be-acquired": {
   "id": "be-acquired",
   "name": "be acquired",
   "event": "being-acquired",
   "nouns": [
    "props",
    "actors"
   ],
   "actions": [
    "aa0f3ad3144d3b009591e22f4f7f9aa2"
   ]
  },
  "be-closed-by": {
   "id": "be-closed-by",
   "name": "be closed by",
   "event": "being-closed-by",
   "nouns": [
    "props",
    "actors"
   ],
   "actions": [
    "cddd2b508635fcd0fcc15232ddfee0d9"
   ]
  },
  "be-discussed": {
   "id": "be-discussed",
   "name": "be discussed",
   "event": "being-discussed",
   "nouns": [
    "quips",
    "actors"
   ],
   "actions": [
    "e8ea14871dd41a378b7d216c974bfdd0"
   ]
  },
  "be-examined": {
   "id": "be-examined",
   "name": "be examined",
   "event": "being-examined",
   "nouns": [
    "objects",
    "actors"
   ],
   "actions": [
    "f5982891e7ebdc570e12599f257d70ea"
   ]
  },
  "be-inserted": {
   "id": "be-inserted",
   "name": "be inserted",
   "event": "being-inserted",
   "nouns": [
    "props",
    "actors",
    "containers"
   ],
   "actions": [
    "e8b94915f1dace4e2b2f59f3a9162c"
   ]
  },
  "be-opened-by": {
   "id": "be-opened-by",
   "name": "be opened by",
   "event": "being-opened-by",
   "nouns": [
    "props",
    "actors"
   ],
   "actions": [
    "cce18d15e14eddefbd4c7741dce6f5ce"
   ]
  },
  "be-passed-through": {
   "id": "be-passed-through",
   "name": "be passed through",
   "event": "being-passed-through",
   "nouns": [
    "doors",
    "actors"
   ],
   "actions": [
    "fe407e25d67b995683a4497a99da6e5e"
   ]
  },
  "close-it": {
   "id": "close-it",
   "name": "close it",
   "event": "closing-it",
   "nouns": [
    "actors",
    "props"
   ],
   "actions": [
    "a028e313823833dd7dad493203bd86fc"
   ]
  },
  "commence": {
   "id": "commence",
   "name": "commence",
   "event": "commencing",
   "nouns": [
    "stories"
   ],
   "actions": [
    "be2fb26f4dfcdf76f7d89ab1a25c254"
   ]
  },
  "comment": {
   "id": "comment",
   "name": "comment",
   "event": "commenting",
   "nouns": [
    "actors",
    "quips"
   ],
   "actions": [
    "efb9c6008baa5f297da22b6b88a69ae"
   ]
  },
  "debug-contents": {
   "id": "debug-contents",
   "name": "debug contents",
   "event": "debugging-contents",
   "nouns": [
    "actors",
    "objects"
   ],
   "actions": [
    "e27d852c08e661f57fed7a37fb44b6"
   ]
  },
  "debug-direct-parent": {
   "id": "debug-direct-parent",
   "name": "debug direct parent",
   "event": "debugging-direct-parent",
   "nouns": [
    "actors",
    "objects"
   ],
   "actions": [
    "bb4ae3b5d6669a4eb5d919a52800816"
   ]
  },
  "debug-room-contents": {
   "id": "debug-room-contents",
   "name": "debug room contents",
   "event": "debugging-room-contents",
   "nouns": [
    "actors"
   ],
   "actions": [
    "efb794467ebba90fb9d54e3f0332632e"
   ]
  },
  "depart": {
   "id": "depart",
   "name": "depart",
   "event": "departing",
   "nouns": [
    "actors"
   ],
   "actions": [
    "ab0570189d883833d088ab9206526dc3"
   ]
  },
  "describe-the-first-room": {
   "id": "describe-the-first-room",
   "name": "describe the first room",
   "event": "describing-the-first-room",
   "nouns": [
    "stories",
    "rooms"
   ],
   "actions": [
    "d4983e203d6eeadd7b989e46d6"
   ]
  },
  "discuss": {
   "id": "discuss",
   "name": "discuss",
   "event": "discussing",
   "nouns": [
    "actors",
    "quips"
   ],
   "actions": [
    "b703927b07424ab1a58cc61fc9287c"
   ]
  },
  "eat-it": {
   "id": "eat-it",
   "name": "eat it",
   "event": "eating-it",
   "nouns": [
    "actors",
    "props"
   ],
   "actions": [
    "dc5f2a70d4f8b2f822f0d4f669c330b"
   ]
  },
  "end-the-story": {
   "id": "end-the-story",
   "name": "end the story",
   "event": "ending-the-story",
   "nouns": [
    "stories"
   ],
   "actions": [
    "a1157c67ba6d0bacc2a189d6ac0f88"
   ]
  },
  "end-turn": {
   "id": "end-turn",
   "name": "end turn",
   "event": "ending-the-turn",
   "nouns": [
    "stories"
   ],
   "actions": [
    "e549441d59279c1801dcd222a57f"
   ]
  },
  "examine-it": {
   "id": "examine-it",
   "name": "examine it",
   "event": "examining-it",
   "nouns": [
    "actors",
    "objects"
   ],
   "actions": [
    "f3d54a0c2d53c10de008fb6ab276670f"
   ]
  },
  "feed-it": {
   "id": "feed-it",
   "name": "feed it",
   "event": "feeding-it",
   "nouns": [
    "actors",
    "objects"
   ],
   "actions": [
    "d87a6fc604d8a034efb8a192199a1a"
   ]
  },
  "give-it-to": {
   "id": "give-it-to",
   "name": "give it to",
   "event": "giving-it-to",
   "nouns": [
    "actors",
    "actors",
    "props"
   ],
   "actions": [
    "c1cd32951100ecf601c99e78a416faa"
   ]
  },
  "go-through-it": {
   "id": "go-through-it",
   "name": "go through it",
   "event": "going-through-it",
   "nouns": [
    "actors",
    "doors"
   ],
   "actions": [
    "f147fff0a14762b6527a8f71240cdf00"
   ]
  },
  "go-to": {
   "id": "go-to",
   "name": "go to",
   "event": "going-to",
   "nouns": [
    "actors",
    "directions"
   ],
   "actions": [
    "a279a3ec2a6f0098be46df9a674f96f"
   ]
  },
  "greet": {
   "id": "greet",
   "name": "greet",
   "event": "greeting",
   "nouns": [
    "actors",
    "actors"
   ],
   "actions": [
    "d73f8083c58056e0b106f0330cabd14e"
   ]
  },
  "impress": {
   "id": "impress",
   "name": "impress",
   "event": "impressing",
   "nouns": [
    "actors"
   ],
   "actions": [
    "a4e06ed8921a1c0e863719f30eed2"
   ]
  },
  "insert-it-into": {
   "id": "insert-it-into",
   "name": "insert it into",
   "event": "inserting-it-into",
   "nouns": [
    "actors",
    "containers",
    "props"
   ],
   "actions": [
    "e7456c6895a9da6aa"
   ]
  },
  "jump": {
   "id": "jump",
   "name": "jump",
   "event": "jumping",
   "nouns": [
    "actors"
   ],
   "actions": [
    "c7e571f4cb23c1766247c3657077f3b9"
   ]
  },
  "kiss-it": {
   "id": "kiss-it",
   "name": "kiss it",
   "event": "kissing-it",
   "nouns": [
    "actors",
    "objects"
   ],
   "actions": [
    "dbc0e2eebefb53eb032ffd9e2de"
   ]
  },
  "listen": {
   "id": "listen",
   "name": "listen",
   "event": "listening",
   "nouns": [
    "actors"
   ],
   "actions": [
    "fce05a152a90e9075438ca19ce7c6"
   ]
  },
  "listen-to": {
   "id": "listen-to",
   "name": "listen to",
   "event": "listening-to",
   "nouns": [
    "actors",
    "kinds"
   ],
   "actions": [
    "fa673c29457c2a9ac47c49658685"
   ]
  },
  "look": {
   "id": "look",
   "name": "look",
   "event": "looking",
   "nouns": [
    "actors"
   ],
   "actions": [
    "dba636750406cd8d7bd0da3"
   ]
  },
  "look-under-it": {
   "id": "look-under-it",
   "name": "look under it",
   "event": "looking-under-it",
   "nouns": [
    "actors",
    "objects"
   ],
   "actions": [
    "b0d267fe6bec1f08710b46623fe815b"
   ]
  },
  "open-it": {
   "id": "open-it",
   "name": "open it",
   "event": "opening-it",
   "nouns": [
    "actors",
    "props"
   ],
   "actions": [
    "b6e604472836e917e8ad1dd6ba2b93"
   ]
  },
  "parse-player-input": {
   "id": "parse-player-input",
   "name": "parse player input",
   "event": "parsing-player-input",
   "nouns": [
    "stories"
   ],
   "actions": null
  },
  "print-contents": {
   "id": "print-contents",
   "name": "print contents",
   "event": "printing-contents",
   "nouns": [
    "objects"
   ],
   "actions": null
  },
  "print-conversation-choices": {
   "id": "print-conversation-choices",
   "name": "print conversation choices",
   "event": "printing-conversation-choices",
   "nouns": [
    "actors",
    "actors"
   ],
   "actions": [
    "e2d83166d1ece71900065119ffe42599"
   ]
  },
  "print-description": {
   "id": "print-description",
   "name": "print description",
   "event": "describing",
   "nouns": [
    "objects"
   ],
   "actions": [
    "cabe59e7088cb62f7dceba9b137394"
   ]
  },
  "print-details": {
   "id": "print-details",
   "name": "print details",
   "event": "printing-details",
   "nouns": [
    "objects"
   ],
   "actions": [
    "e5d3985209c83582a8cec2119c5f808"
   ]
  },
  "print-name": {
   "id": "print-name",
   "name": "print name",
   "event": "printing-name-text",
   "nouns": [
    "objects"
   ],
   "actions": [
    "f3a039df8a36cf28033c549d1d"
   ]
  },
  "print-the-banner": {
   "id": "print-the-banner",
   "name": "print the banner",
   "event": "printing-the-banner",
   "nouns": [
    "stories"
   ],
   "actions": [
    "f80bca4a7a24cc9ee42067363ecda08"
   ]
  },
  "put-it-onto": {
   "id": "put-it-onto",
   "name": "put it onto",
   "event": "putting-it-onto",
   "nouns": [
    "actors",
    "supporters",
    "props"
   ],
   "actions": [
    "c7ebd1166c4885183f433a6de51411e1"
   ]
  },
  "receive-insertion": {
   "id": "receive-insertion",
   "name": "receive insertion",
   "event": "receiving-insertion",
   "nouns": [
    "containers",
    "props",
    "actors"
   ],
   "actions": [
    "bd9ceebcdcc7ebb57be7187da0b3916"
   ]
  },
  "report-already-closed": {
   "id": "report-already-closed",
   "name": "report already closed",
   "event": "reporting-already-closed",
   "nouns": [
    "props",
    "actors"
   ],
   "actions": [
    "a0cf32d8daa137026c62f8494129"
   ]
  },
  "report-already-off": {
   "id": "report-already-off",
   "name": "report already off",
   "event": "reporting-already-off",
   "nouns": [
    "devices",
    "actors"
   ],
   "actions": [
    "e54bc28e4a6d6ee8f15cde43d041f41"
   ]
  },
  "report-already-on": {
   "id": "report-already-on",
   "name": "report already on",
   "event": "reporting-already-on",
   "nouns": [
    "devices",
    "actors"
   ],
   "actions": [
    "b44ad3f16c77942ddebfbcbbdb1381c"
   ]
  },
  "report-already-open": {
   "id": "report-already-open",
   "name": "report already open",
   "event": "reporting-already-opened",
   "nouns": [
    "props",
    "actors"
   ],
   "actions": [
    "f15ec508566495ef319bda0247c3ba"
   ]
  },
  "report-attack": {
   "id": "report-attack",
   "name": "report attack",
   "event": "reporting-attack",
   "nouns": [
    "objects",
    "actors"
   ],
   "actions": [
    "e092832c33702a363a16fa5d40421bf"
   ]
  },
  "report-comment": {
   "id": "report-comment",
   "name": "report comment",
   "event": "reporting-comment",
   "nouns": [
    "quips",
    "actors"
   ],
   "actions": [
    "f1487317dd37aafc7fdc5f81a5296e88"
   ]
  },
  "report-currently-closed": {
   "id": "report-currently-closed",
   "name": "report currently closed",
   "event": "reporting-currently-closed",
   "nouns": [
    "doors",
    "actors"
   ],
   "actions": [
    "fdb1dc013d4b6490c6586654cf3"
   ]
  },
  "report-eat": {
   "id": "report-eat",
   "name": "report eat",
   "event": "reporting-eat",
   "nouns": [
    "props",
    "actors"
   ],
   "actions": [
    "f4d2a3f2bc4ff9069889a9299b"
   ]
  },
  "report-gave": {
   "id": "report-gave",
   "name": "report gave",
   "event": "reporting-gave",
   "nouns": [
    "props",
    "actors",
    "actors"
   ],
   "actions": [
    "f820558419c997630a2044117eacd4b9"
   ]
  },
  "report-give": {
   "id": "report-give",
   "name": "report give",
   "event": "reporting-give",
   "nouns": [
    "actors",
    "props",
    "actors"
   ],
   "actions": [
    "bb5aa921e7bf8a1a2a243bc496da3d37"
   ]
  },
  "report-inoperable": {
   "id": "report-inoperable",
   "name": "report inoperable",
   "event": "reporting-inoperable",
   "nouns": [
    "props"
   ],
   "actions": [
    "fdcb93cc0324bd63a08f6415d5f62e3"
   ]
  },
  "report-inventory": {
   "id": "report-inventory",
   "name": "report inventory",
   "event": "reporting-inventory",
   "nouns": [
    "actors"
   ],
   "actions": [
    "fa0bacce6442c0445e80e1fea52eb2"
   ]
  },
  "report-jump": {
   "id": "report-jump",
   "name": "report jump",
   "event": "reporting-jump",
   "nouns": [
    "kinds",
    "actors"
   ],
   "actions": [
    "fd97981f521bc27a06faf20b0d5053da"
   ]
  },
  "report-kiss": {
   "id": "report-kiss",
   "name": "report kiss",
   "event": "reporting-kiss",
   "nouns": [
    "objects",
    "actors"
   ],
   "actions": [
    "a03c941d25050d3a9974d9165bb8f1"
   ]
  },
  "report-listen": {
   "id": "report-listen",
   "name": "report listen",
   "event": "reporting-listen",
   "nouns": [
    "kinds",
    "actors"
   ],
   "actions": [
    "b0b8e208c4cc475e100e14fbcf2f4ae"
   ]
  },
  "report-locked": {
   "id": "report-locked",
   "name": "report locked",
   "event": "reporting-locked",
   "nouns": [
    "props",
    "actors"
   ],
   "actions": [
    "b53ef45a3e1238a3aa47a2bdc8d9da4"
   ]
  },
  "report-look-under": {
   "id": "report-look-under",
   "name": "report look under",
   "event": "reporting-look-under",
   "nouns": [
    "objects",
    "actors"
   ],
   "actions": [
    "eecdc7b5555289906c2b415fb201cb9b"
   ]
  },
  "report-not-closeable": {
   "id": "report-not-closeable",
   "name": "report not closeable",
   "event": "reporting-not-closeable",
   "nouns": [
    "props",
    "actors"
   ],
   "actions": [
    "d621a08c3de87b0ee4cfe7971a9"
   ]
  },
  "report-now-closed": {
   "id": "report-now-closed",
   "name": "report now closed",
   "event": "reporting-now-closed",
   "nouns": [
    "props",
    "actors"
   ],
   "actions": [
    "a14a3d5d1cca145f13f50ca1309af311"
   ]
  },
  "report-now-off": {
   "id": "report-now-off",
   "name": "report now off",
   "event": "reporting-now-off",
   "nouns": [
    "devices",
    "actors"
   ],
   "actions": [
    "bc33e94a52fffabe6b20bcc54955ce3"
   ]
  },
  "report-now-on": {
   "id": "report-now-on",
   "name": "report now on",
   "event": "reporting-now-on",
   "nouns": [
    "devices",
    "actors"
   ],
   "actions": [
    "f13522698fb406b905407ccdfec39"
   ]
  },
  "report-now-open": {
   "id": "report-now-open",
   "name": "report now open",
   "event": "reporting-now-open",
   "nouns": [
    "props",
    "actors"
   ],
   "actions": [
    "cd40d14902a558771bcc65444bc871f"
   ]
  },
  "report-placed": {
   "id": "report-placed",
   "name": "report placed",
   "event": "reporting-placed",
   "nouns": [
    "props",
    "actors",
    "supporters"
   ],
   "actions": [
    "a0ca28fd1e1f6003219b799c59e5925"
   ]
  },
  "report-put": {
   "id": "report-put",
   "name": "report put",
   "event": "reporting-put",
   "nouns": [
    "supporters",
    "props",
    "actors"
   ],
   "actions": [
    "eba487ee4a161c0ed3a629161b612e"
   ]
  },
  "report-search": {
   "id": "report-search",
   "name": "report search",
   "event": "reporting-search",
   "nouns": [
    "props",
    "actors"
   ],
   "actions": [
    "c346a8302cd8e56c361e834a576fde7"
   ]
  },
  "report-show": {
   "id": "report-show",
   "name": "report show",
   "event": "reporting-show",
   "nouns": [
    "actors",
    "props",
    "actors"
   ],
   "actions": [
    "ea0713e6af256dfe532c004ff5a23f5"
   ]
  },
  "report-shown": {
   "id": "report-shown",
   "name": "report shown",
   "event": "reporting-shown",
   "nouns": [
    "props",
    "actors",
    "actors"
   ],
   "actions": [
    "e460183c17f770fec808d132736f24e"
   ]
  },
  "report-smell": {
   "id": "report-smell",
   "name": "report smell",
   "event": "reporting-smell",
   "nouns": [
    "kinds",
    "actors"
   ],
   "actions": [
    "fe0b565e4a5208b37419400dbbaea13"
   ]
  },
  "report-switch-off": {
   "id": "report-switch-off",
   "name": "report switch off",
   "event": "reporting-switch-off",
   "nouns": [
    "devices"
   ],
   "actions": [
    "ef7b600f685b87fcd10d15cb21a3"
   ]
  },
  "report-switched-on": {
   "id": "report-switched-on",
   "name": "report switched on",
   "event": "reporting-switched-on",
   "nouns": [
    "devices",
    "actors"
   ],
   "actions": [
    "b703969791fe6a4b4b43764fe9215df0"
   ]
  },
  "report-take": {
   "id": "report-take",
   "name": "report take",
   "event": "reporting-take",
   "nouns": [
    "props",
    "actors"
   ],
   "actions": [
    "c9d7e290147df7748f2222edd327fcb2"
   ]
  },
  "report-the-view": {
   "id": "report-the-view",
   "name": "report the view",
   "event": "reporting-the-view",
   "nouns": [
    "rooms"
   ],
   "actions": [
    "be3c2951941a98acc45dfd2d4705308"
   ]
  },
  "report-unopenable": {
   "id": "report-unopenable",
   "name": "report unopenable",
   "event": "reporting-unopenable",
   "nouns": [
    "props",
    "actors"
   ],
   "actions": [
    "c76bb26a42aa9719a9a0313b2ec1193d"
   ]
  },
  "report-wear": {
   "id": "report-wear",
   "name": "report wear",
   "event": "reporting-wear",
   "nouns": [
    "props",
    "actors"
   ],
   "actions": [
    "dba3f28c029c8ebf74fe7ac2d94e0ac0"
   ]
  },
  "search-it": {
   "id": "search-it",
   "name": "search it",
   "event": "searching-it",
   "nouns": [
    "actors",
    "props"
   ],
   "actions": [
    "ed97f386ce027a2ff676ef7a8a302e1"
   ]
  },
  "set-initial-position": {
   "id": "set-initial-position",
   "name": "set initial position",
   "event": "setting-initial-position",
   "nouns": [
    "stories",
    "actors",
    "rooms"
   ],
   "actions": [
    "f6b153d8dfbb91bec15a99f842d2"
   ]
  },
  "show-it-to": {
   "id": "show-it-to",
   "name": "show it to",
   "event": "showing-it-to",
   "nouns": [
    "actors",
    "actors",
    "props"
   ],
   "actions": [
    "b5bd6ab4d4df9616f49967630082b1e"
   ]
  },
  "smell": {
   "id": "smell",
   "name": "smell",
   "event": "smelling",
   "nouns": [
    "actors"
   ],
   "actions": [
    "bd1fabe270be6ab6a4724cc8f142da1"
   ]
  },
  "smell-it": {
   "id": "smell-it",
   "name": "smell it",
   "event": "smelling-it",
   "nouns": [
    "actors",
    "kinds"
   ],
   "actions": [
    "f1eb81a12f8e22e8c9964742f7a703"
   ]
  },
  "switch-it-off": {
   "id": "switch-it-off",
   "name": "switch it off",
   "event": "switching-it-off",
   "nouns": [
    "actors",
    "props"
   ],
   "actions": [
    "d61637c94b5914263b9b2e888c6ee3f0"
   ]
  },
  "switch-it-on": {
   "id": "switch-it-on",
   "name": "switch it on",
   "event": "switching-it-on",
   "nouns": [
    "actors",
    "props"
   ],
   "actions": [
    "c8f1338cac6671fc5d6920c2c0f7b893"
   ]
  },
  "take-it": {
   "id": "take-it",
   "name": "take it",
   "event": "taking-it",
   "nouns": [
    "actors",
    "props"
   ],
   "actions": [
    "ad645ba09d5207de6d3ab358466b7355"
   ]
  },
  "wear-it": {
   "id": "wear-it",
   "name": "wear it",
   "event": "wearing-it",
   "nouns": [
    "actors",
    "props"
   ],
   "actions": [
    "a2afcc3585516a895061346c3d5fa3f"
   ]
  }
 },
 "classes": {
  "actors": {
   "id": "actors",
   "parents": [
    "objects",
    "kinds"
   ],
   "plural": "actors",
   "singular": "actor",
   "properties": [
    {
     "id": "actors-not-inputing-dialog",
     "type": "EnumProperty",
     "name": "not inputing dialog"
    },
    {
     "id": "actors-clothing",
     "type": "PointerProperty",
     "name": "clothing",
     "relates": "objects",
     "relation": "actors-clothing-relation",
     "many": true
    },
    {
     "id": "actors-inventory",
     "type": "PointerProperty",
     "name": "inventory",
     "relates": "objects",
     "relation": "actors-inventory-relation",
     "many": true
    },
    {
     "id": "actors-greeting",
     "type": "PointerProperty",
     "name": "greeting",
     "relates": "quips"
    },
    {
     "id": "actors-next-quip",
     "type": "PointerProperty",
     "name": "next quip",
     "relates": "quips"
    }
   ]
  },
  "animals": {
   "id": "animals",
   "parents": [
    "actors",
    "objects",
    "kinds"
   ],
   "plural": "animals",
   "singular": "animal",
   "properties": null
  },
  "canisters": {
   "id": "canisters",
   "parents": [
    "containers",
    "openers",
    "props",
    "objects",
    "kinds"
   ],
   "plural": "canisters",
   "singular": "canister",
   "properties": [
    {
     "id": "canisters-hidden",
     "type": "EnumProperty",
     "name": "hidden"
    }
   ]
  },
  "containers": {
   "id": "containers",
   "parents": [
    "openers",
    "props",
    "objects",
    "kinds"
   ],
   "plural": "containers",
   "singular": "container",
   "properties": [
    {
     "id": "containers-opaque",
     "type": "EnumProperty",
     "name": "opaque"
    },
    {
     "id": "containers-lockable",
     "type": "EnumProperty",
     "name": "lockable"
    },
    {
     "id": "containers-locked",
     "type": "EnumProperty",
     "name": "locked"
    },
    {
     "id": "containers-contents",
     "type": "PointerProperty",
     "name": "contents",
     "relates": "objects",
     "relation": "containers-contents-relation",
     "many": true
    }
   ]
  },
  "conversation-globals": {
   "id": "conversation-globals",
   "parents": [
    "kinds"
   ],
   "plural": "conversation globals",
   "singular": "conversation global",
   "properties": [
    {
     "id": "conversation-globals-queue",
     "type": "PointerProperty",
     "name": "queue",
     "relates": "quips",
     "many": true
    },
    {
     "id": "conversation-globals-interlocutor",
     "type": "PointerProperty",
     "name": "interlocutor",
     "relates": "actors"
    },
    {
     "id": "conversation-globals-parent",
     "type": "PointerProperty",
     "name": "parent",
     "relates": "quips"
    },
    {
     "id": "conversation-globals-grandparent",
     "type": "PointerProperty",
     "name": "grandparent",
     "relates": "quips"
    },
    {
     "id": "conversation-globals-greatgrand",
     "type": "PointerProperty",
     "name": "greatgrand",
     "relates": "quips"
    }
   ]
  },
  "devices": {
   "id": "devices",
   "parents": [
    "props",
    "objects",
    "kinds"
   ],
   "plural": "devices",
   "singular": "device",
   "properties": [
    {
     "id": "devices-operable",
     "type": "EnumProperty",
     "name": "operable"
    },
    {
     "id": "devices-switched-off",
     "type": "EnumProperty",
     "name": "switched off"
    }
   ]
  },
  "directions": {
   "id": "directions",
   "parents": [
    "kinds"
   ],
   "plural": "directions",
   "singular": "direction",
   "properties": [
    {
     "id": "directions-opposite",
     "type": "PointerProperty",
     "name": "opposite",
     "relates": "directions",
     "relation": "directions-opposite-relation"
    },
    {
     "id": "directions-x-opposite",
     "type": "PointerProperty",
     "name": "x-opposite",
     "relates": "directions",
     "relation": "directions-opposite-relation"
    }
   ]
  },
  "doors": {
   "id": "doors",
   "parents": [
    "openers",
    "props",
    "objects",
    "kinds"
   ],
   "plural": "doors",
   "singular": "door",
   "properties": [
    {
     "id": "doors-x-via-down",
     "type": "PointerProperty",
     "name": "x-via-down",
     "relates": "rooms",
     "relation": "rooms-down-via-relation",
     "many": true
    },
    {
     "id": "doors-destination",
     "type": "PointerProperty",
     "name": "destination",
     "relates": "doors",
     "relation": "doors-destination-relation"
    },
    {
     "id": "doors-sources",
     "type": "PointerProperty",
     "name": "sources",
     "relates": "doors",
     "relation": "doors-destination-relation",
     "many": true
    },
    {
     "id": "doors-x-via-north",
     "type": "PointerProperty",
     "name": "x-via-north",
     "relates": "rooms",
     "relation": "rooms-north-via-relation",
     "many": true
    },
    {
     "id": "doors-x-via-south",
     "type": "PointerProperty",
     "name": "x-via-south",
     "relates": "rooms",
     "relation": "rooms-south-via-relation",
     "many": true
    },
    {
     "id": "doors-x-via-east",
     "type": "PointerProperty",
     "name": "x-via-east",
     "relates": "rooms",
     "relation": "rooms-east-via-relation",
     "many": true
    },
    {
     "id": "doors-x-via-west",
     "type": "PointerProperty",
     "name": "x-via-west",
     "relates": "rooms",
     "relation": "rooms-west-via-relation",
     "many": true
    },
    {
     "id": "doors-x-via-up",
     "type": "PointerProperty",
     "name": "x-via-up",
     "relates": "rooms",
     "relation": "rooms-up-via-relation",
     "many": true
    }
   ]
  },
  "facts": {
   "id": "facts",
   "parents": [
    "kinds"
   ],
   "plural": "facts",
   "singular": "fact",
   "properties": [
    {
     "id": "facts-summary",
     "type": "TextProperty",
     "name": "summary"
    }
   ]
  },
  "following-quips": {
   "id": "following-quips",
   "parents": [
    "kinds"
   ],
   "plural": "following quips",
   "singular": "following quip",
   "properties": [
    {
     "id": "following-quips-following",
     "type": "PointerProperty",
     "name": "following",
     "relates": "quips"
    },
    {
     "id": "following-quips-indirectly-following",
     "type": "EnumProperty",
     "name": "indirectly following"
    },
    {
     "id": "following-quips-leading",
     "type": "PointerProperty",
     "name": "leading",
     "relates": "quips"
    }
   ]
  },
  "kinds": {
   "id": "kinds",
   "parents": null,
   "plural": "kinds",
   "singular": "kind",
   "properties": [
    {
     "id": "kinds-printed-name",
     "type": "TextProperty",
     "name": "printed name"
    },
    {
     "id": "kinds-indefinite-article",
     "type": "TextProperty",
     "name": "indefinite article"
    },
    {
     "id": "kinds-singular-named",
     "type": "EnumProperty",
     "name": "singular-named"
    },
    {
     "id": "kinds-common-named",
     "type": "EnumProperty",
     "name": "common-named"
    },
    {
     "id": "kinds-recollected",
     "type": "EnumProperty",
     "name": "recollected"
    },
    {
     "id": "kinds-name",
     "type": "TextProperty",
     "name": "name"
    }
   ]
  },
  "objects": {
   "id": "objects",
   "parents": [
    "kinds"
   ],
   "plural": "objects",
   "singular": "object",
   "properties": [
    {
     "id": "objects-wearer",
     "type": "PointerProperty",
     "name": "wearer",
     "relates": "actors",
     "relation": "actors-clothing-relation"
    },
    {
     "id": "objects-unhandled",
     "type": "EnumProperty",
     "name": "unhandled"
    },
    {
     "id": "objects-enclosure",
     "type": "PointerProperty",
     "name": "enclosure",
     "relates": "containers",
     "relation": "containers-contents-relation"
    },
    {
     "id": "objects-brief",
     "type": "TextProperty",
     "name": "brief"
    },
    {
     "id": "objects-scenery",
     "type": "EnumProperty",
     "name": "scenery"
    },
    {
     "id": "objects-whereabouts",
     "type": "PointerProperty",
     "name": "whereabouts",
     "relates": "rooms",
     "relation": "rooms-contents-relation"
    },
    {
     "id": "objects-owner",
     "type": "PointerProperty",
     "name": "owner",
     "relates": "actors",
     "relation": "actors-inventory-relation"
    },
    {
     "id": "objects-support",
     "type": "PointerProperty",
     "name": "support",
     "relates": "supporters",
     "relation": "supporters-contents-relation"
    },
    {
     "id": "objects-description",
     "type": "TextProperty",
     "name": "description"
    }
   ]
  },
  "openers": {
   "id": "openers",
   "parents": [
    "props",
    "objects",
    "kinds"
   ],
   "plural": "openers",
   "singular": "opener",
   "properties": [
    {
     "id": "openers-unlocked",
     "type": "EnumProperty",
     "name": "unlocked"
    },
    {
     "id": "openers-open",
     "type": "EnumProperty",
     "name": "open"
    },
    {
     "id": "openers-hinged",
     "type": "EnumProperty",
     "name": "hinged"
    },
    {
     "id": "openers-locakable",
     "type": "EnumProperty",
     "name": "locakable"
    }
   ]
  },
  "pending-quips": {
   "id": "pending-quips",
   "parents": [
    "kinds"
   ],
   "plural": "pending quips",
   "singular": "pending quip",
   "properties": [
    {
     "id": "pending-quips-subject",
     "type": "PointerProperty",
     "name": "subject",
     "relates": "actors"
    },
    {
     "id": "pending-quips-immediate",
     "type": "EnumProperty",
     "name": "immediate"
    },
    {
     "id": "pending-quips-obligatory",
     "type": "EnumProperty",
     "name": "obligatory"
    }
   ]
  },
  "props": {
   "id": "props",
   "parents": [
    "objects",
    "kinds"
   ],
   "plural": "props",
   "singular": "prop",
   "properties": [
    {
     "id": "props-wearable",
     "type": "EnumProperty",
     "name": "wearable"
    },
    {
     "id": "props-portable",
     "type": "EnumProperty",
     "name": "portable"
    }
   ]
  },
  "quip-requirements": {
   "id": "quip-requirements",
   "parents": [
    "kinds"
   ],
   "plural": "quip requirements",
   "singular": "quip requirement",
   "properties": [
    {
     "id": "quip-requirements-fact",
     "type": "PointerProperty",
     "name": "fact",
     "relates": "facts"
    },
    {
     "id": "quip-requirements-quip",
     "type": "PointerProperty",
     "name": "quip",
     "relates": "quips"
    },
    {
     "id": "quip-requirements-permitted",
     "type": "EnumProperty",
     "name": "permitted"
    }
   ]
  },
  "quips": {
   "id": "quips",
   "parents": [
    "kinds"
   ],
   "plural": "quips",
   "singular": "quip",
   "properties": [
    {
     "id": "quips-planned",
     "type": "EnumProperty",
     "name": "planned"
    },
    {
     "id": "quips-comment",
     "type": "TextProperty",
     "name": "comment"
    },
    {
     "id": "quips-subject",
     "type": "PointerProperty",
     "name": "subject",
     "relates": "actors"
    },
    {
     "id": "quips-reply",
     "type": "TextProperty",
     "name": "reply"
    },
    {
     "id": "quips-repeatable",
     "type": "EnumProperty",
     "name": "repeatable"
    },
    {
     "id": "quips-restrictive",
     "type": "EnumProperty",
     "name": "restrictive"
    }
   ]
  },
  "rooms": {
   "id": "rooms",
   "parents": [
    "kinds"
   ],
   "plural": "rooms",
   "singular": "room",
   "properties": [
    {
     "id": "rooms-description",
     "type": "TextProperty",
     "name": "description"
    },
    {
     "id": "rooms-west-via",
     "type": "PointerProperty",
     "name": "west-via",
     "relates": "doors",
     "relation": "rooms-west-via-relation"
    },
    {
     "id": "rooms-down-via",
     "type": "PointerProperty",
     "name": "down-via",
     "relates": "doors",
     "relation": "rooms-down-via-relation"
    },
    {
     "id": "rooms-south-via",
     "type": "PointerProperty",
     "name": "south-via",
     "relates": "doors",
     "relation": "rooms-south-via-relation"
    },
    {
     "id": "rooms-visited",
     "type": "EnumProperty",
     "name": "visited"
    },
    {
     "id": "rooms-contents",
     "type": "PointerProperty",
     "name": "contents",
     "relates": "objects",
     "relation": "rooms-contents-relation",
     "many": true
    },
    {
     "id": "rooms-north-via",
     "type": "PointerProperty",
     "name": "north-via",
     "relates": "doors",
     "relation": "rooms-north-via-relation"
    },
    {
     "id": "rooms-east-via",
     "type": "PointerProperty",
     "name": "east-via",
     "relates": "doors",
     "relation": "rooms-east-via-relation"
    },
    {
     "id": "rooms-up-via",
     "type": "PointerProperty",
     "name": "up-via",
     "relates": "doors",
     "relation": "rooms-up-via-relation"
    }
   ]
  },
  "status-bar-instances": {
   "id": "status-bar-instances",
   "parents": [
    "kinds"
   ],
   "plural": "status bar instances",
   "singular": "status bar instance",
   "properties": [
    {
     "id": "status-bar-instances-left",
     "type": "TextProperty",
     "name": "left"
    },
    {
     "id": "status-bar-instances-right",
     "type": "TextProperty",
     "name": "right"
    }
   ]
  },
  "stories": {
   "id": "stories",
   "parents": [
    "kinds"
   ],
   "plural": "stories",
   "singular": "story",
   "properties": [
    {
     "id": "stories-player-input",
     "type": "TextProperty",
     "name": "player input"
    },
    {
     "id": "stories-scored",
     "type": "EnumProperty",
     "name": "scored"
    },
    {
     "id": "stories-playing",
     "type": "EnumProperty",
     "name": "playing"
    },
    {
     "id": "stories-author",
     "type": "TextProperty",
     "name": "author"
    },
    {
     "id": "stories-headline",
     "type": "TextProperty",
     "name": "headline"
    },
    {
     "id": "stories-score",
     "type": "NumProperty",
     "name": "score"
    },
    {
     "id": "stories-maximum-score",
     "type": "NumProperty",
     "name": "maximum score"
    },
    {
     "id": "stories-turn-count",
     "type": "NumProperty",
     "name": "turn count"
    }
   ]
  },
  "supporters": {
   "id": "supporters",
   "parents": [
    "props",
    "objects",
    "kinds"
   ],
   "plural": "supporters",
   "singular": "supporter",
   "properties": [
    {
     "id": "supporters-contents",
     "type": "PointerProperty",
     "name": "contents",
     "relates": "objects",
     "relation": "supporters-contents-relation",
     "many": true
    }
   ]
  }
 },
 "enums": {
  "actors-not-inputing-dialog": {
   "choices": [
    "not-inputing-dialog",
    "inputing-dialog"
   ]
  },
  "canisters-hidden": {
   "choices": [
    "hidden",
    "found"
   ]
  },
  "containers-lockable": {
   "choices": [
    "not-lockable",
    "lockable"
   ]
  },
  "containers-locked": {
   "choices": [
    "unlocked",
    "locked"
   ]
  },
  "containers-opaque": {
   "choices": [
    "opaque",
    "transparent"
   ]
  },
  "devices-operable": {
   "choices": [
    "operable",
    "inoperable"
   ]
  },
  "devices-switched-off": {
   "choices": [
    "switched-off",
    "switched-on"
   ]
  },
  "following-quips-indirectly-following": {
   "choices": [
    "indirectly-following",
    "directly-following"
   ]
  },
  "kinds-common-named": {
   "choices": [
    "common-named",
    "proper-named"
   ]
  },
  "kinds-recollected": {
   "choices": [
    "not-recollected",
    "recollected"
   ]
  },
  "kinds-singular-named": {
   "choices": [
    "singular-named",
    "plural-named"
   ]
  },
  "objects-scenery": {
   "choices": [
    "not-scenery",
    "scenery"
   ]
  },
  "objects-unhandled": {
   "choices": [
    "unhandled",
    "handled"
   ]
  },
  "openers-hinged": {
   "choices": [
    "hinged",
    "hingeless"
   ]
  },
  "openers-locakable": {
   "choices": [
    "not-lockable",
    "locakable"
   ]
  },
  "openers-open": {
   "choices": [
    "open",
    "closed"
   ]
  },
  "openers-unlocked": {
   "choices": [
    "unlocked",
    "locked"
   ]
  },
  "pending-quips-immediate": {
   "choices": [
    "immediate",
    "postponed"
   ]
  },
  "pending-quips-obligatory": {
   "choices": [
    "obligatory",
    "optional"
   ]
  },
  "props-portable": {
   "choices": [
    "portable",
    "fixed-in-place"
   ]
  },
  "props-wearable": {
   "choices": [
    "not-wearable",
    "wearable"
   ]
  },
  "quip-requirements-permitted": {
   "choices": [
    "permitted",
    "prohibited"
   ]
  },
  "quips-planned": {
   "choices": [
    "planned",
    "casual"
   ]
  },
  "quips-repeatable": {
   "choices": [
    "one-time",
    "repeatable"
   ]
  },
  "quips-restrictive": {
   "choices": [
    "unrestricted",
    "restrictive"
   ]
  },
  "rooms-visited": {
   "choices": [
    "unvisited",
    "visited"
   ]
  },
  "stories-playing": {
   "choices": [
    "playing",
    "completed"
   ]
  },
  "stories-scored": {
   "choices": [
    "unscored",
    "scored"
   ]
  }
 },
 "events": {
  "acquiring-it": {
   "id": "acquiring-it",
   "name": "acquiring it"
  },
  "attacking-it": {
   "id": "attacking-it",
   "name": "attacking it"
  },
  "being-acquired": {
   "id": "being-acquired",
   "name": "being acquired"
  },
  "being-closed-by": {
   "id": "being-closed-by",
   "name": "being closed by"
  },
  "being-discussed": {
   "id": "being-discussed",
   "name": "being discussed"
  },
  "being-examined": {
   "id": "being-examined",
   "name": "being examined",
   "capture": [
    {
     "Instance": "paints",
     "Class": "props",
     "Callback": "b0b3278650e296e6c97cd7ab4ad833a",
     "Options": 3
    },
    {
     "Instance": "cloths",
     "Class": "props",
     "Callback": "aa170312949d078407f725833b2f5060",
     "Options": 3
    },
    {
     "Instance": "painting",
     "Class": "props",
     "Callback": "bb18d25b98660857d88e5f8d3747cdd",
     "Options": 1
    },
    {
     "Instance": "painting",
     "Class": "props",
     "Callback": "f5a8b9fe48a05ba88d82829dc5e83",
     "Options": 3
    },
    {
     "Instance": "table",
     "Class": "supporters",
     "Callback": "f7de12be85b890ff8a18141f91c5c976",
     "Options": 3
    },
    {
     "Instance": "telegram",
     "Class": "props",
     "Callback": "a27a04f5b30dcdb2fb3bca460ffa83a3",
     "Options": 3
    },
    {
     "Instance": "gravel",
     "Class": "props",
     "Callback": "f18bf4192989ca36ad6671e40b253546",
     "Options": 3
    },
    {
     "Instance": "seaweed",
     "Class": "props",
     "Callback": "d01cb53d42b7e1ff4ae8bf8f8c468188",
     "Options": 3
    },
    {
     "Instance": "evil-fish",
     "Class": "animals",
     "Callback": "afe1549f4f3744702939c3f69275690",
     "Options": 3
    },
    {
     "Instance": "bouquet",
     "Class": "props",
     "Callback": "deb7efe1826c5a9b5819f6073fbc089a",
     "Options": 3
    },
    {
     "Instance": "lingerie-bag",
     "Class": "props",
     "Callback": "f6f9049ec302b58175ad90021534b902",
     "Options": 3
    }
   ],
   "bubble": [
    {
     "Instance": "aquarium",
     "Class": "containers",
     "Callback": "e9cbdbcdcc4ceac8f98fc99adca7268",
     "Options": 1
    }
   ]
  },
  "being-inserted": {
   "id": "being-inserted",
   "name": "being inserted",
   "bubble": [
    {
     "Instance": "fish-food",
     "Class": "canisters",
     "Callback": "f383a1073f44b75c602b1e7b404e",
     "Options": 1
    }
   ]
  },
  "being-opened-by": {
   "id": "being-opened-by",
   "name": "being opened by",
   "bubble": [
    {
     "Instance": "fish-food",
     "Class": "canisters",
     "Callback": "c3955aa0d0fc455e328e1faa99333f88",
     "Options": 1
    }
   ]
  },
  "being-passed-through": {
   "id": "being-passed-through",
   "name": "being passed through"
  },
  "closing-it": {
   "id": "closing-it",
   "name": "closing it"
  },
  "commencing": {
   "id": "commencing",
   "name": "commencing",
   "bubble": [
    {
     "Instance": "day-for-fresh-sushi",
     "Class": "stories",
     "Callback": "f386909950f2484579711506a8baff3f",
     "Options": 1
    }
   ]
  },
  "commenting": {
   "id": "commenting",
   "name": "commenting"
  },
  "debugging-contents": {
   "id": "debugging-contents",
   "name": "debugging contents"
  },
  "debugging-direct-parent": {
   "id": "debugging-direct-parent",
   "name": "debugging direct parent"
  },
  "debugging-room-contents": {
   "id": "debugging-room-contents",
   "name": "debugging room contents"
  },
  "departing": {
   "id": "departing",
   "name": "departing"
  },
  "describing": {
   "id": "describing",
   "name": "describing"
  },
  "describing-the-first-room": {
   "id": "describing-the-first-room",
   "name": "describing the first room"
  },
  "discussing": {
   "id": "discussing",
   "name": "discussing"
  },
  "eating-it": {
   "id": "eating-it",
   "name": "eating it"
  },
  "ending-the-story": {
   "id": "ending-the-story",
   "name": "ending the story"
  },
  "ending-the-turn": {
   "id": "ending-the-turn",
   "name": "ending the turn",
   "capture": [
    {
     "Instance": "",
     "Class": "stories",
     "Callback": "e994180f6a5d75a96cbdd199f4385b7",
     "Options": 1
    }
   ],
   "bubble": [
    {
     "Instance": "",
     "Class": "stories",
     "Callback": "c937462c89e5bdcf0dd8d0922d98",
     "Options": 1
    },
    {
     "Instance": "",
     "Class": "stories",
     "Callback": "d3d7b2c9b0c8ddbc95756b981d73e5f",
     "Options": 1
    }
   ]
  },
  "examining-it": {
   "id": "examining-it",
   "name": "examining it"
  },
  "feeding-it": {
   "id": "feeding-it",
   "name": "feeding it",
   "capture": [
    {
     "Instance": "",
     "Class": "actors",
     "Callback": "cfa04e62556a30555426ce11",
     "Options": 1
    }
   ],
   "bubble": [
    {
     "Instance": "",
     "Class": "actors",
     "Callback": "ec0f84674c17bbc66902dfe2b2ab284",
     "Options": 1
    },
    {
     "Instance": "",
     "Class": "actors",
     "Callback": "b35f0712de0ecce8238feaa309229a9",
     "Options": 1
    }
   ]
  },
  "giving-it-to": {
   "id": "giving-it-to",
   "name": "giving it to",
   "capture": [
    {
     "Instance": "",
     "Class": "actors",
     "Callback": "be81d6d6fb1282452bb02cb39edd0ab7",
     "Options": 0
    },
    {
     "Instance": "",
     "Class": "actors",
     "Callback": "bbbc9f864e3bc39e4f95286879cfb14",
     "Options": 0
    },
    {
     "Instance": "",
     "Class": "actors",
     "Callback": "d81e1bcae38ca5faa6da34c0abc814e",
     "Options": 0
    }
   ]
  },
  "going-through-it": {
   "id": "going-through-it",
   "name": "going through it"
  },
  "going-to": {
   "id": "going-to",
   "name": "going to"
  },
  "greeting": {
   "id": "greeting",
   "name": "greeting"
  },
  "impressing": {
   "id": "impressing",
   "name": "impressing"
  },
  "inserting-it-into": {
   "id": "inserting-it-into",
   "name": "inserting it into",
   "capture": [
    {
     "Instance": "",
     "Class": "actors",
     "Callback": "a174d13dbae209ddf968f9f16e9",
     "Options": 0
    },
    {
     "Instance": "",
     "Class": "actors",
     "Callback": "e74361077dff040e0f5ac2c",
     "Options": 0
    },
    {
     "Instance": "",
     "Class": "actors",
     "Callback": "c092ea1a525fdfe12f7b255a354",
     "Options": 0
    }
   ]
  },
  "jumping": {
   "id": "jumping",
   "name": "jumping",
   "bubble": [
    {
     "Instance": "player",
     "Class": "actors",
     "Callback": "d3aad960bedd830691389807b6897",
     "Options": 1
    }
   ]
  },
  "kissing-it": {
   "id": "kissing-it",
   "name": "kissing it",
   "capture": [
    {
     "Instance": "",
     "Class": "actors",
     "Callback": "f9baafd0c928bb023c1ae00466c0464f",
     "Options": 0
    }
   ]
  },
  "listening": {
   "id": "listening",
   "name": "listening"
  },
  "listening-to": {
   "id": "listening-to",
   "name": "listening to"
  },
  "looking": {
   "id": "looking",
   "name": "looking"
  },
  "looking-under-it": {
   "id": "looking-under-it",
   "name": "looking under it"
  },
  "opening-it": {
   "id": "opening-it",
   "name": "opening it"
  },
  "parsing-player-input": {
   "id": "parsing-player-input",
   "name": "parsing player input",
   "bubble": [
    {
     "Instance": "",
     "Class": "stories",
     "Callback": "c5f3dd965e3297200df816e37587df",
     "Options": 1
    }
   ]
  },
  "printing-contents": {
   "id": "printing-contents",
   "name": "printing contents",
   "bubble": [
    {
     "Instance": "",
     "Class": "containers",
     "Callback": "b8fe8f66db86e05e87b099ee578f0c",
     "Options": 1
    },
    {
     "Instance": "",
     "Class": "supporters",
     "Callback": "dcddc49d44fe886a8f96d952f",
     "Options": 1
    }
   ]
  },
  "printing-conversation-choices": {
   "id": "printing-conversation-choices",
   "name": "printing conversation choices"
  },
  "printing-details": {
   "id": "printing-details",
   "name": "printing details",
   "bubble": [
    {
     "Instance": "studio",
     "Class": "rooms",
     "Callback": "d95e07eb79eaf3fa2f3302d5ec0bf09",
     "Options": 1
    },
    {
     "Instance": "window",
     "Class": "openers",
     "Callback": "a5c72be52c4f2336436967baf5754621",
     "Options": 1
    }
   ]
  },
  "printing-name-text": {
   "id": "printing-name-text",
   "name": "printing name text",
   "bubble": [
    {
     "Instance": "",
     "Class": "containers",
     "Callback": "ec75951720bb443a92e928bf812c3",
     "Options": 1
    },
    {
     "Instance": "",
     "Class": "doors",
     "Callback": "fbe785a6bfa575e678ea90da0211546",
     "Options": 1
    },
    {
     "Instance": "",
     "Class": "devices",
     "Callback": "d5a222bae6a821c254abbe62ff2d879",
     "Options": 1
    }
   ]
  },
  "printing-the-banner": {
   "id": "printing-the-banner",
   "name": "printing the banner"
  },
  "putting-it-onto": {
   "id": "putting-it-onto",
   "name": "putting it onto",
   "capture": [
    {
     "Instance": "",
     "Class": "actors",
     "Callback": "fe5e3c26bd18f5891ae038584f8d",
     "Options": 0
    },
    {
     "Instance": "",
     "Class": "actors",
     "Callback": "ce6c9fbff93abf1dbaa2997c357e86dc",
     "Options": 0
    },
    {
     "Instance": "",
     "Class": "actors",
     "Callback": "a8e0802cdf9350defff5d67b08b6",
     "Options": 0
    },
    {
     "Instance": "",
     "Class": "actors",
     "Callback": "adc65b6c27ecd7b4aabf6bdd973b",
     "Options": 0
    }
   ]
  },
  "receiving-insertion": {
   "id": "receiving-insertion",
   "name": "receiving insertion",
   "capture": [
    {
     "Instance": "",
     "Class": "containers",
     "Callback": "ff8957b596f3ee3ccd61c07445b3e2c8",
     "Options": 0
    }
   ],
   "bubble": [
    {
     "Instance": "vase",
     "Class": "containers",
     "Callback": "f12f24b77e1c05dfffb8569d4a8cb088",
     "Options": 1
    }
   ]
  },
  "reporting-already-closed": {
   "id": "reporting-already-closed",
   "name": "reporting already closed"
  },
  "reporting-already-off": {
   "id": "reporting-already-off",
   "name": "reporting already off"
  },
  "reporting-already-on": {
   "id": "reporting-already-on",
   "name": "reporting already on"
  },
  "reporting-already-opened": {
   "id": "reporting-already-opened",
   "name": "reporting already opened"
  },
  "reporting-attack": {
   "id": "reporting-attack",
   "name": "reporting attack",
   "bubble": [
    {
     "Instance": "evil-fish",
     "Class": "animals",
     "Callback": "dfc7b1de160d5dea8bc04e39ff06899c",
     "Options": 1
    }
   ]
  },
  "reporting-comment": {
   "id": "reporting-comment",
   "name": "reporting comment"
  },
  "reporting-currently-closed": {
   "id": "reporting-currently-closed",
   "name": "reporting currently closed"
  },
  "reporting-eat": {
   "id": "reporting-eat",
   "name": "reporting eat"
  },
  "reporting-gave": {
   "id": "reporting-gave",
   "name": "reporting gave",
   "bubble": [
    {
     "Instance": "fish-food",
     "Class": "canisters",
     "Callback": "b9f230f5cf3e2fed1850f455e7",
     "Options": 1
    }
   ]
  },
  "reporting-give": {
   "id": "reporting-give",
   "name": "reporting give"
  },
  "reporting-inoperable": {
   "id": "reporting-inoperable",
   "name": "reporting inoperable"
  },
  "reporting-inventory": {
   "id": "reporting-inventory",
   "name": "reporting inventory"
  },
  "reporting-jump": {
   "id": "reporting-jump",
   "name": "reporting jump"
  },
  "reporting-kiss": {
   "id": "reporting-kiss",
   "name": "reporting kiss",
   "bubble": [
    {
     "Instance": "evil-fish",
     "Class": "animals",
     "Callback": "a5003a18e42eccfb92e40558889d715",
     "Options": 1
    }
   ]
  },
  "reporting-listen": {
   "id": "reporting-listen",
   "name": "reporting listen"
  },
  "reporting-locked": {
   "id": "reporting-locked",
   "name": "reporting locked"
  },
  "reporting-look-under": {
   "id": "reporting-look-under",
   "name": "reporting look under",
   "bubble": [
    {
     "Instance": "cabinet",
     "Class": "containers",
     "Callback": "ca750763ef503a6b3d7b11b970801455",
     "Options": 1
    },
    {
     "Instance": "cloths",
     "Class": "props",
     "Callback": "c5351efbb7da045f6635e48a3e1958d",
     "Options": 1
    },
    {
     "Instance": "table",
     "Class": "supporters",
     "Callback": "cf141e4b8a35226dfdb1ea917e278a",
     "Options": 1
    }
   ]
  },
  "reporting-not-closeable": {
   "id": "reporting-not-closeable",
   "name": "reporting not closeable"
  },
  "reporting-now-closed": {
   "id": "reporting-now-closed",
   "name": "reporting now closed",
   "bubble": [
    {
     "Instance": "cabinet",
     "Class": "containers",
     "Callback": "df1fbacb62e40599f5b43f7e52eb",
     "Options": 1
    }
   ]
  },
  "reporting-now-off": {
   "id": "reporting-now-off",
   "name": "reporting now off"
  },
  "reporting-now-on": {
   "id": "reporting-now-on",
   "name": "reporting now on"
  },
  "reporting-now-open": {
   "id": "reporting-now-open",
   "name": "reporting now open",
   "bubble": [
    {
     "Instance": "cabinet",
     "Class": "containers",
     "Callback": "c239f2d147803170152e8566a7",
     "Options": 1
    },
    {
     "Instance": "window",
     "Class": "openers",
     "Callback": "a60e331213831e69ad2e4bc2294f1",
     "Options": 1
    }
   ]
  },
  "reporting-placed": {
   "id": "reporting-placed",
   "name": "reporting placed"
  },
  "reporting-put": {
   "id": "reporting-put",
   "name": "reporting put"
  },
  "reporting-search": {
   "id": "reporting-search",
   "name": "reporting search",
   "bubble": [
    {
     "Instance": "cloths",
     "Class": "props",
     "Callback": "c5351efbb7da045f6635e48a3e1958d",
     "Options": 1
    }
   ]
  },
  "reporting-show": {
   "id": "reporting-show",
   "name": "reporting show"
  },
  "reporting-shown": {
   "id": "reporting-shown",
   "name": "reporting shown",
   "bubble": [
    {
     "Instance": "cloths",
     "Class": "props",
     "Callback": "a795e203a59b3f39adb4b28ca8eef4c3",
     "Options": 1
    }
   ]
  },
  "reporting-smell": {
   "id": "reporting-smell",
   "name": "reporting smell",
   "bubble": [
    {
     "Instance": "studio",
     "Class": "rooms",
     "Callback": "a98fd9211e9759d489fb4522c66d852f",
     "Options": 1
    },
    {
     "Instance": "bouquet",
     "Class": "props",
     "Callback": "ea5f6c77ffaf10b74b109e760",
     "Options": 1
    }
   ]
  },
  "reporting-switch-off": {
   "id": "reporting-switch-off",
   "name": "reporting switch off"
  },
  "reporting-switched-on": {
   "id": "reporting-switched-on",
   "name": "reporting switched on"
  },
  "reporting-take": {
   "id": "reporting-take",
   "name": "reporting take",
   "capture": [
    {
     "Instance": "",
     "Class": "doors",
     "Callback": "d6fb6189d5ad739d1f229b238d88c",
     "Options": 1
    }
   ],
   "bubble": [
    {
     "Instance": "paints",
     "Class": "props",
     "Callback": "d87833ece051f1acbf636ab531b13708",
     "Options": 1
    },
    {
     "Instance": "painting",
     "Class": "props",
     "Callback": "de68b28196dad365db8afd062982",
     "Options": 1
    },
    {
     "Instance": "evil-fish",
     "Class": "animals",
     "Callback": "e81a1f6706715fbda9c4a073ca7ef98",
     "Options": 1
    }
   ]
  },
  "reporting-the-view": {
   "id": "reporting-the-view",
   "name": "reporting the view"
  },
  "reporting-unopenable": {
   "id": "reporting-unopenable",
   "name": "reporting unopenable"
  },
  "reporting-wear": {
   "id": "reporting-wear",
   "name": "reporting wear"
  },
  "searching-it": {
   "id": "searching-it",
   "name": "searching it"
  },
  "setting-initial-position": {
   "id": "setting-initial-position",
   "name": "setting initial position"
  },
  "showing-it-to": {
   "id": "showing-it-to",
   "name": "showing it to",
   "capture": [
    {
     "Instance": "",
     "Class": "actors",
     "Callback": "aea3c76840abada678c70a2122fe83ff",
     "Options": 0
    },
    {
     "Instance": "",
     "Class": "actors",
     "Callback": "c2f1521949bf2b49274bdf1aacd83349",
     "Options": 0
    }
   ]
  },
  "smelling": {
   "id": "smelling",
   "name": "smelling"
  },
  "smelling-it": {
   "id": "smelling-it",
   "name": "smelling it"
  },
  "switching-it-off": {
   "id": "switching-it-off",
   "name": "switching it off"
  },
  "switching-it-on": {
   "id": "switching-it-on",
   "name": "switching it on"
  },
  "taking-it": {
   "id": "taking-it",
   "name": "taking it"
  },
  "wearing-it": {
   "id": "wearing-it",
   "name": "wearing it"
  }
 },
 "instances": {
  "aquarium": {
   "id": "aquarium",
   "type": "containers",
   "name": "aquarium",
   "values": {
    "containers-opaque": "transparent",
    "kinds-name": "aquarium",
    "objects-brief": "In one corner of the room, a large aquarium bubbles in menacing fashion.",
    "objects-description": "A very roomy aquarium, large enough to hold quite a variety of colorful sealife -- if any yet survived.",
    "objects-enclosure": "",
    "objects-owner": "",
    "objects-support": "",
    "objects-wearer": "",
    "objects-whereabouts": "studio",
    "openers-hinged": "hingeless",
    "openers-open": "open",
    "props-portable": "fixed-in-place"
   }
  },
  "bouquet": {
   "id": "bouquet",
   "type": "props",
   "name": "bouquet",
   "values": {
    "kinds-name": "bouquet",
    "objects-description": "Okay, so it's silly and sentimental and no doubt a waste of money, of which there is never really enough, but: you miss her. You've missed her since ten seconds after she stepped aboard the shuttle to Luna Prime, and when you saw these -- her favorites, pure golden tulips like springtime -- you had to have them.",
    "objects-enclosure": "",
    "objects-owner": "player",
    "objects-support": "",
    "objects-wearer": "",
    "objects-whereabouts": ""
   }
  },
  "britney": {
   "id": "britney",
   "type": "actors",
   "name": "Britney",
   "values": {
    "kinds-name": "Britney",
    "objects-enclosure": "",
    "objects-owner": "",
    "objects-support": "",
    "objects-wearer": "",
    "objects-whereabouts": ""
   }
  },
  "cabinet": {
   "id": "cabinet",
   "type": "containers",
   "name": "cabinet",
   "values": {
    "kinds-name": "cabinet",
    "objects-brief": "A huge cabinet, in the guise of an armoire, stands between the windows.",
    "objects-description": "Large, and with a bit of an Art Nouveau theme going on in the shape of the doors.",
    "objects-enclosure": "",
    "objects-owner": "",
    "objects-support": "",
    "objects-wearer": "",
    "objects-whereabouts": "studio",
    "openers-hinged": "hinged",
    "openers-open": "closed",
    "props-portable": "fixed-in-place"
   }
  },
  "chef-hat": {
   "id": "chef-hat",
   "type": "props",
   "name": "chef hat",
   "values": {
    "kinds-name": "chef hat",
    "objects-description": "A big white chef hat of the kind worn by chefs. In this case, you. Just goes to show what a hurry you were in on the way out of the restaurant.",
    "objects-enclosure": "",
    "objects-owner": "",
    "objects-support": "",
    "objects-wearer": "player",
    "objects-whereabouts": ""
   }
  },
  "closed-cabinet": {
   "id": "closed-cabinet",
   "type": "facts",
   "name": "closedCabinet",
   "values": {
    "kinds-name": "closedCabinet"
   }
  },
  "cloths": {
   "id": "cloths",
   "type": "props",
   "name": "cloths",
   "values": {
    "kinds-name": "cloths",
    "kinds-singular-named": "plural-named",
    "objects-description": "Various colors of drapery that Britney uses to set up backgrounds and clothe her models. She does a lot of portraiture, so this comes in handy. It's all a big messy wad at the moment. Organized is not her middle name.",
    "objects-enclosure": "cabinet",
    "objects-owner": "",
    "objects-support": "",
    "objects-wearer": "",
    "objects-whereabouts": ""
   }
  },
  "conversation": {
   "id": "conversation",
   "type": "conversation-globals",
   "name": "conversation",
   "values": {
    "kinds-name": "conversation"
   }
  },
  "day-for-fresh-sushi": {
   "id": "day-for-fresh-sushi",
   "type": "stories",
   "name": "Day For Fresh Sushi",
   "values": {
    "kinds-name": "Day For Fresh Sushi",
    "stories-author": "Emily Short",
    "stories-headline": "Your basic surreal gay fish romance",
    "stories-maximum-score": 1,
    "stories-scored": "scored"
   }
  },
  "down": {
   "id": "down",
   "type": "directions",
   "name": "down",
   "values": {
    "directions-opposite": "up",
    "directions-x-opposite": "up",
    "kinds-name": "down"
   }
  },
  "easel": {
   "id": "easel",
   "type": "supporters",
   "name": "easel",
   "values": {
    "kinds-name": "easel",
    "objects-enclosure": "",
    "objects-owner": "",
    "objects-scenery": "scenery",
    "objects-support": "",
    "objects-wearer": "",
    "objects-whereabouts": "studio"
   }
  },
  "east": {
   "id": "east",
   "type": "directions",
   "name": "east",
   "values": {
    "directions-opposite": "west",
    "directions-x-opposite": "west",
    "kinds-name": "east"
   }
  },
  "evil-fish": {
   "id": "evil-fish",
   "type": "animals",
   "name": "evil fish",
   "values": {
    "kinds-name": "evil fish",
    "objects-description": "Even if you had had no prior experience with him, you would be able to see at a glance that this is an evil fish. From his sharkish nose to his razor fins, every inch of his compact body exudes hatred and danger.",
    "objects-enclosure": "aquarium",
    "objects-owner": "",
    "objects-support": "",
    "objects-wearer": "",
    "objects-whereabouts": ""
   }
  },
  "examined-bag-once": {
   "id": "examined-bag-once",
   "type": "facts",
   "name": "examinedBagOnce",
   "values": {
    "kinds-name": "examinedBagOnce"
   }
  },
  "examined-bag-twice": {
   "id": "examined-bag-twice",
   "type": "facts",
   "name": "examinedBagTwice",
   "values": {
    "kinds-name": "examinedBagTwice"
   }
  },
  "examined-bouquet": {
   "id": "examined-bouquet",
   "type": "facts",
   "name": "examinedBouquet",
   "values": {
    "kinds-name": "examinedBouquet"
   }
  },
  "examined-cloths": {
   "id": "examined-cloths",
   "type": "facts",
   "name": "examinedCloths",
   "values": {
    "kinds-name": "examinedCloths"
   }
  },
  "examined-fish-once": {
   "id": "examined-fish-once",
   "type": "facts",
   "name": "examinedFishOnce",
   "values": {
    "kinds-name": "examinedFishOnce"
   }
  },
  "examined-fish-twice": {
   "id": "examined-fish-twice",
   "type": "facts",
   "name": "examinedFishTwice",
   "values": {
    "kinds-name": "examinedFishTwice"
   }
  },
  "examined-gravel": {
   "id": "examined-gravel",
   "type": "facts",
   "name": "examinedGravel",
   "values": {
    "kinds-name": "examinedGravel"
   }
  },
  "examined-painting": {
   "id": "examined-painting",
   "type": "facts",
   "name": "examinedPainting",
   "values": {
    "kinds-name": "examinedPainting"
   }
  },
  "examined-paints": {
   "id": "examined-paints",
   "type": "facts",
   "name": "examinedPaints",
   "values": {
    "kinds-name": "examinedPaints"
   }
  },
  "examined-seaweed": {
   "id": "examined-seaweed",
   "type": "facts",
   "name": "examinedSeaweed",
   "values": {
    "kinds-name": "examinedSeaweed"
   }
  },
  "examined-telegraph": {
   "id": "examined-telegraph",
   "type": "facts",
   "name": "examinedTelegraph",
   "values": {
    "kinds-name": "examinedTelegraph"
   }
  },
  "fish-food": {
   "id": "fish-food",
   "type": "canisters",
   "name": "fish food",
   "values": {
    "canisters-hidden": "hidden",
    "containers-opaque": "opaque",
    "kinds-name": "fish food",
    "kinds-singular-named": "plural-named",
    "objects-description": "A vehemently orange canister of fish food.",
    "objects-enclosure": "",
    "objects-owner": "",
    "objects-support": "",
    "objects-wearer": "",
    "objects-whereabouts": "",
    "openers-hinged": "hingeless",
    "openers-open": "closed"
   }
  },
  "gravel": {
   "id": "gravel",
   "type": "props",
   "name": "gravel",
   "values": {
    "kinds-name": "gravel",
    "kinds-singular-named": "plural-named",
    "objects-description": "A lot of very small grey rocks.",
    "objects-enclosure": "aquarium",
    "objects-owner": "",
    "objects-support": "",
    "objects-wearer": "",
    "objects-whereabouts": ""
   }
  },
  "inserted-flowers": {
   "id": "inserted-flowers",
   "type": "facts",
   "name": "insertedFlowers",
   "values": {
    "kinds-name": "insertedFlowers"
   }
  },
  "lingerie-bag": {
   "id": "lingerie-bag",
   "type": "props",
   "name": "lingerie bag",
   "values": {
    "kinds-name": "lingerie bag",
    "objects-description": "You grant yourself the satisfaction of a little peek inside. You went with a pale, silky ivory this time -- it has that kind of sophisticated innocence, and it goes well with the purple of your skin. A small smirk of anticipation crosses your lips.",
    "objects-enclosure": "",
    "objects-owner": "player",
    "objects-support": "",
    "objects-wearer": "",
    "objects-whereabouts": ""
   }
  },
  "looked-under-cabinet": {
   "id": "looked-under-cabinet",
   "type": "facts",
   "name": "lookedUnderCabinet",
   "values": {
    "kinds-name": "lookedUnderCabinet"
   }
  },
  "looked-under-table": {
   "id": "looked-under-table",
   "type": "facts",
   "name": "lookedUnderTable",
   "values": {
    "kinds-name": "lookedUnderTable"
   }
  },
  "north": {
   "id": "north",
   "type": "directions",
   "name": "north",
   "values": {
    "directions-opposite": "south",
    "directions-x-opposite": "south",
    "kinds-name": "north"
   }
  },
  "opened-cabinet": {
   "id": "opened-cabinet",
   "type": "facts",
   "name": "openedCabinet",
   "values": {
    "kinds-name": "openedCabinet"
   }
  },
  "opened-window": {
   "id": "opened-window",
   "type": "facts",
   "name": "openedWindow",
   "values": {
    "kinds-name": "openedWindow"
   }
  },
  "painting": {
   "id": "painting",
   "type": "props",
   "name": "painting",
   "values": {
    "kinds-name": "painting",
    "objects-description": "Only partway finished, but you can tell what it is: Britney's mother. You only met the old woman once, before she faded out of existence in a little hospice in Salzburg.\n\nIn the picture, her hands are grasping tightly at a small grey bottle, the pills to which she became addicted in her old age, and strange, gargoyle-like forms clutch at her arms and whisper in her ears.\n\nBut the disturbing thing, the truly awful thing, is the small figure of Britney herself, down in the corner, unmistakable: she is walking away. Her back turned.\n\nYou thought she'd finally talked this out, but evidently not. Still feels guilty for leaving. You only barely stop yourself from tracing, with your finger, those tiny slumped shoulders...\n",
    "objects-enclosure": "",
    "objects-owner": "",
    "objects-support": "easel",
    "objects-wearer": "",
    "objects-whereabouts": ""
   }
  },
  "paints": {
   "id": "paints",
   "type": "props",
   "name": "paints",
   "values": {
    "kinds-name": "paints",
    "kinds-singular-named": "plural-named",
    "objects-description": "A bunch of tubes of oil paint, most of them in some state of grunginess, some with the tops twisted partway off.",
    "objects-enclosure": "cabinet",
    "objects-owner": "",
    "objects-support": "",
    "objects-wearer": "",
    "objects-whereabouts": ""
   }
  },
  "player": {
   "id": "player",
   "type": "actors",
   "name": "player",
   "values": {
    "kinds-name": "player",
    "objects-enclosure": "",
    "objects-owner": "",
    "objects-scenery": "scenery",
    "objects-support": "",
    "objects-wearer": "",
    "objects-whereabouts": "studio"
   }
  },
  "seaweed": {
   "id": "seaweed",
   "type": "props",
   "name": "seaweed",
   "values": {
    "kinds-name": "seaweed",
    "kinds-singular-named": "plural-named",
    "objects-description": "Fake plastic seaweed of the kind generally bought in stores for exactly this purpose.",
    "objects-enclosure": "aquarium",
    "objects-owner": "",
    "objects-support": "",
    "objects-wearer": "",
    "objects-whereabouts": ""
   }
  },
  "smelled-bouquet": {
   "id": "smelled-bouquet",
   "type": "facts",
   "name": "smelledBouquet",
   "values": {
    "kinds-name": "smelledBouquet"
   }
  },
  "south": {
   "id": "south",
   "type": "directions",
   "name": "south",
   "values": {
    "directions-opposite": "north",
    "directions-x-opposite": "north",
    "kinds-name": "south"
   }
  },
  "status-bar": {
   "id": "status-bar",
   "type": "status-bar-instances",
   "name": "status bar",
   "values": {
    "kinds-name": "status bar"
   }
  },
  "studio": {
   "id": "studio",
   "type": "rooms",
   "name": "studio",
   "values": {
    "kinds-name": "Studio",
    "kinds-printed-name": "Studio",
    "rooms-down-via": "",
    "rooms-east-via": "",
    "rooms-north-via": "",
    "rooms-south-via": "",
    "rooms-up-via": "",
    "rooms-west-via": ""
   }
  },
  "table": {
   "id": "table",
   "type": "supporters",
   "name": "table",
   "values": {
    "kinds-name": "table",
    "objects-description": "A monstrosity of poor taste and bad design: made of some heavy, French-empire sort of wood, with a single pillar for a central leg, carved in the image of Poseidon surrounded by nymphs. It's all scaley, and whenever you sit down, the trident has a tendency to stab you in the knee. But Britney assures you it's worth a fortune.",
    "objects-enclosure": "",
    "objects-owner": "",
    "objects-scenery": "scenery",
    "objects-support": "",
    "objects-wearer": "",
    "objects-whereabouts": "studio"
   }
  },
  "telegram": {
   "id": "telegram",
   "type": "props",
   "name": "telegram",
   "values": {
    "kinds-name": "telegram",
    "objects-description": "A telegram, apparently. And dated three days ago. TRIUMPH OURS STOP BACK SOON STOP BE SURE TO FEED FISH STOP",
    "objects-enclosure": "",
    "objects-owner": "player",
    "objects-support": "",
    "objects-wearer": "",
    "objects-whereabouts": ""
   }
  },
  "took-paints": {
   "id": "took-paints",
   "type": "facts",
   "name": "tookPaints",
   "values": {
    "kinds-name": "tookPaints"
   }
  },
  "up": {
   "id": "up",
   "type": "directions",
   "name": "up",
   "values": {
    "directions-opposite": "down",
    "directions-x-opposite": "down",
    "kinds-name": "up"
   }
  },
  "vase": {
   "id": "vase",
   "type": "containers",
   "name": "vase",
   "values": {
    "kinds-name": "vase",
    "objects-description": "A huge vase -- what you saw once described in a Regency romance as an epergne, maybe -- something so big that it would block someone sitting at the table from seeing anyone else also sitting at the table. But it does function nicely as a receptacle for hugeass bouquets of flowers.",
    "objects-enclosure": "",
    "objects-owner": "",
    "objects-support": "table",
    "objects-wearer": "",
    "objects-whereabouts": "",
    "openers-hinged": "hingeless",
    "openers-open": "open"
   }
  },
  "west": {
   "id": "west",
   "type": "directions",
   "name": "west",
   "values": {
    "directions-opposite": "east",
    "directions-x-opposite": "east",
    "kinds-name": "west"
   }
  },
  "window": {
   "id": "window",
   "type": "openers",
   "name": "window",
   "values": {
    "kinds-name": "window",
    "objects-enclosure": "",
    "objects-owner": "",
    "objects-scenery": "scenery",
    "objects-support": "",
    "objects-wearer": "",
    "objects-whereabouts": "studio",
    "openers-hinged": "hinged",
    "openers-open": "closed"
   }
  }
 },
 "aliases": {
  "aquarium": [
   "aquarium",
   "aquarium"
  ],
  "armoire": [
   "cabinet"
  ],
  "bag": [
   "examined-bag-twice",
   "examined-bag-once",
   "lingerie-bag"
  ],
  "bar": [
   "status-bar"
  ],
  "big": [
   "chef-hat"
  ],
  "bouquet": [
   "bouquet",
   "bouquet",
   "examined-bouquet",
   "smelled-bouquet"
  ],
  "britney": [
   "britney",
   "britney"
  ],
  "cabinet": [
   "cabinet",
   "closed-cabinet",
   "opened-cabinet",
   "cabinet",
   "looked-under-cabinet"
  ],
  "can": [
   "fish-food"
  ],
  "chef": [
   "chef-hat"
  ],
  "chef hat": [
   "chef-hat"
  ],
  "chef's": [
   "chef-hat"
  ],
  "chefs": [
   "chef-hat"
  ],
  "closed": [
   "closed-cabinet"
  ],
  "closed cabinet": [
   "closed-cabinet"
  ],
  "cloth": [
   "cloths"
  ],
  "cloths": [
   "cloths",
   "examined-cloths",
   "cloths"
  ],
  "conversation": [
   "conversation",
   "conversation"
  ],
  "d": [
   "down"
  ],
  "day": [
   "day-for-fresh-sushi"
  ],
  "day for fresh sushi": [
   "day-for-fresh-sushi"
  ],
  "down": [
   "down",
   "down"
  ],
  "drapery": [
   "cloths"
  ],
  "e": [
   "east"
  ],
  "easel": [
   "easel",
   "easel"
  ],
  "east": [
   "east",
   "east"
  ],
  "evil": [
   "evil-fish"
  ],
  "evil fish": [
   "evil-fish"
  ],
  "examined": [
   "examined-gravel",
   "examined-bag-twice",
   "examined-bag-once",
   "examined-cloths",
   "examined-paints",
   "examined-fish-once",
   "examined-bouquet",
   "examined-telegraph",
   "examined-fish-twice",
   "examined-seaweed",
   "examined-painting"
  ],
  "examined bag once": [
   "examined-bag-once"
  ],
  "examined bag twice": [
   "examined-bag-twice"
  ],
  "examined bouquet": [
   "examined-bouquet"
  ],
  "examined cloths": [
   "examined-cloths"
  ],
  "examined fish once": [
   "examined-fish-once"
  ],
  "examined fish twice": [
   "examined-fish-twice"
  ],
  "examined gravel": [
   "examined-gravel"
  ],
  "examined painting": [
   "examined-painting"
  ],
  "examined paints": [
   "examined-paints"
  ],
  "examined seaweed": [
   "examined-seaweed"
  ],
  "examined telegraph": [
   "examined-telegraph"
  ],
  "fish": [
   "evil-fish",
   "evil-fish",
   "fish-food",
   "examined-fish-once",
   "examined-fish-twice"
  ],
  "fish food": [
   "fish-food"
  ],
  "flowers": [
   "bouquet",
   "inserted-flowers"
  ],
  "food": [
   "fish-food"
  ],
  "for": [
   "day-for-fresh-sushi"
  ],
  "fresh": [
   "day-for-fresh-sushi"
  ],
  "gravel": [
   "gravel",
   "gravel",
   "examined-gravel"
  ],
  "hat": [
   "chef-hat"
  ],
  "image": [
   "painting"
  ],
  "inserted": [
   "inserted-flowers"
  ],
  "inserted flowers": [
   "inserted-flowers"
  ],
  "lingerie": [
   "lingerie-bag"
  ],
  "lingerie bag": [
   "lingerie-bag"
  ],
  "little rocks": [
   "gravel"
  ],
  "looked": [
   "looked-under-cabinet",
   "looked-under-table"
  ],
  "looked under cabinet": [
   "looked-under-cabinet"
  ],
  "looked under table": [
   "looked-under-table"
  ],
  "n": [
   "north"
  ],
  "north": [
   "north",
   "north"
  ],
  "once": [
   "examined-bag-once",
   "examined-fish-once"
  ],
  "opened": [
   "opened-window",
   "opened-cabinet"
  ],
  "opened cabinet": [
   "opened-cabinet"
  ],
  "opened window": [
   "opened-window"
  ],
  "painting": [
   "painting",
   "examined-painting",
   "painting"
  ],
  "paints": [
   "paints",
   "took-paints",
   "paints",
   "examined-paints"
  ],
  "player": [
   "player",
   "player"
  ],
  "portrait": [
   "painting"
  ],
  "s": [
   "south"
  ],
  "seaweed": [
   "seaweed",
   "seaweed",
   "examined-seaweed"
  ],
  "smelled": [
   "smelled-bouquet"
  ],
  "smelled bouquet": [
   "smelled-bouquet"
  ],
  "south": [
   "south",
   "south"
  ],
  "status": [
   "status-bar"
  ],
  "status bar": [
   "status-bar"
  ],
  "studio": [
   "studio",
   "studio"
  ],
  "sushi": [
   "day-for-fresh-sushi"
  ],
  "table": [
   "table",
   "table",
   "looked-under-table"
  ],
  "tank": [
   "aquarium"
  ],
  "telegram": [
   "telegram",
   "telegram"
  ],
  "telegraph": [
   "examined-telegraph"
  ],
  "took": [
   "took-paints"
  ],
  "took paints": [
   "took-paints"
  ],
  "tulip": [
   "bouquet"
  ],
  "tulips": [
   "bouquet"
  ],
  "twice": [
   "examined-bag-twice",
   "examined-fish-twice"
  ],
  "u": [
   "up"
  ],
  "under": [
   "looked-under-cabinet",
   "looked-under-table"
  ],
  "up": [
   "up",
   "up"
  ],
  "vase": [
   "vase",
   "vase"
  ],
  "w": [
   "west"
  ],
  "weed": [
   "seaweed"
  ],
  "west": [
   "west",
   "west"
  ],
  "white": [
   "chef-hat"
  ],
  "window": [
   "window",
   "window",
   "opened-window"
  ],
  "windows": [
   "window"
  ],
  "yellow paper": [
   "telegram"
  ]
 },
 "parsings": [
  {
   "Action": "attack-it",
   "Commands": [
    "attack|break|smash|hit|fight|torture {{something}}",
    "wreck|crack|destroy|murder|kill|punch|thump {{something}}"
   ]
  },
  {
   "Action": "report-inventory",
   "Commands": [
    "inventory|inv|i"
   ]
  },
  {
   "Action": "look",
   "Commands": [
    "look|l"
   ]
  },
  {
   "Action": "look-under-it",
   "Commands": [
    "look under {{something}}"
   ]
  },
  {
   "Action": "debug-direct-parent",
   "Commands": [
    "parent of {{something}}"
   ]
  },
  {
   "Action": "debug-contents",
   "Commands": [
    "contents of {{something}}"
   ]
  },
  {
   "Action": "debug-room-contents",
   "Commands": [
    "contents of room"
   ]
  },
  {
   "Action": "switch-it-on",
   "Commands": [
    "switch|turn on {{something}}",
    "switch {{something}} on"
   ]
  },
  {
   "Action": "switch-it-off",
   "Commands": [
    "turn|switch off {{something}}",
    "turn|switch {{something}} off"
   ]
  },
  {
   "Action": "eat-it",
   "Commands": [
    "eat {{something}}"
   ]
  },
  {
   "Action": "examine-it",
   "Commands": [
    "examine|x|watch|describe|check {{something}}",
    "look|l {{something}}",
    "look|l at {{something}}"
   ]
  },
  {
   "Action": "give-it-to",
   "Commands": [
    "give|pay|offer|feed {{something}} {{something else}}",
    "give|pay|offer|feed {{something else}} to {{something}}"
   ]
  },
  {
   "Action": "insert-it-into",
   "Commands": [
    "put|insert {{something else}} in|inside|into {{something}}",
    "drop {{something else}} in|into|down {{something}}"
   ]
  },
  {
   "Action": "jump",
   "Commands": [
    "jump|skip|hop"
   ]
  },
  {
   "Action": "kiss-it",
   "Commands": [
    "kiss|hug|embrace {{something}}"
   ]
  },
  {
   "Action": "listen-to",
   "Commands": [
    "listen to {{something}}",
    "listen {{something}}"
   ]
  },
  {
   "Action": "listen",
   "Commands": [
    "listen"
   ]
  },
  {
   "Action": "go-to",
   "Commands": [
    "go {{something}}"
   ]
  },
  {
   "Action": "go-through-it",
   "Commands": [
    "enter {{something}}"
   ]
  },
  {
   "Action": "open-it",
   "Commands": [
    "open {{something}}"
   ]
  },
  {
   "Action": "close-it",
   "Commands": [
    "close {{something}}"
   ]
  },
  {
   "Action": "put-it-onto",
   "Commands": [
    "put {{something else}} on|onto {{something}}",
    "drop|throw|discard {{something else}} on|onto {{something}}"
   ]
  },
  {
   "Action": "search-it",
   "Commands": [
    "search {{something}}",
    "look inside|in|into|through {{something}}"
   ]
  },
  {
   "Action": "show-it-to",
   "Commands": [
    "show|present|display {{something}} {{something else}}",
    "show|present|display {{something else}} to {{something}}"
   ]
  },
  {
   "Action": "smell-it",
   "Commands": [
    "smell|sniff {{something}}"
   ]
  },
  {
   "Action": "smell",
   "Commands": [
    "smell|sniff"
   ]
  },
  {
   "Action": "take-it",
   "Commands": [
    "take|get {{something}}",
    "pick up {{something}}",
    "pick {{something}} up"
   ]
  },
  {
   "Action": "wear-it",
   "Commands": [
    "wear|don {{something}}",
    "put on {{something}}",
    "put {{something}} on"
   ]
  },
  {
   "Action": "greet",
   "Commands": [
    "talk to {{something}}",
    "t|talk|greet|ask {{something}}"
   ]
  },
  {
   "Action": "feed-it",
   "Commands": [
    "feed {{something}}"
   ]
  }
 ],
 "relation": {
  "actors-clothing-relation": {
   "id": "actors-clothing-relation",
   "name": "actors-clothing-relation",
   "style": "OneToMany",
   "src": "actors-clothing",
   "tgt": "objects-wearer"
  },
  "actors-inventory-relation": {
   "id": "actors-inventory-relation",
   "name": "actors-inventory-relation",
   "style": "OneToMany",
   "src": "actors-inventory",
   "tgt": "objects-owner"
  },
  "containers-contents-relation": {
   "id": "containers-contents-relation",
   "name": "containers-contents-relation",
   "style": "OneToMany",
   "src": "containers-contents",
   "tgt": "objects-enclosure"
  },
  "directions-opposite-relation": {
   "id": "directions-opposite-relation",
   "name": "directions-opposite-relation",
   "style": "OneToOne",
   "src": "directions-opposite",
   "tgt": "directions-x-opposite"
  },
  "doors-destination-relation": {
   "id": "doors-destination-relation",
   "name": "doors-destination-relation",
   "style": "ManyToOne",
   "src": "doors-destination",
   "tgt": "doors-sources"
  },
  "rooms-contents-relation": {
   "id": "rooms-contents-relation",
   "name": "rooms-contents-relation",
   "style": "OneToMany",
   "src": "rooms-contents",
   "tgt": "objects-whereabouts"
  },
  "rooms-down-via-relation": {
   "id": "rooms-down-via-relation",
   "name": "rooms-down-via-relation",
   "style": "ManyToOne",
   "src": "rooms-down-via",
   "tgt": "doors-x-via-down"
  },
  "rooms-east-via-relation": {
   "id": "rooms-east-via-relation",
   "name": "rooms-east-via-relation",
   "style": "ManyToOne",
   "src": "rooms-east-via",
   "tgt": "doors-x-via-east"
  },
  "rooms-north-via-relation": {
   "id": "rooms-north-via-relation",
   "name": "rooms-north-via-relation",
   "style": "ManyToOne",
   "src": "rooms-north-via",
   "tgt": "doors-x-via-north"
  },
  "rooms-south-via-relation": {
   "id": "rooms-south-via-relation",
   "name": "rooms-south-via-relation",
   "style": "ManyToOne",
   "src": "rooms-south-via",
   "tgt": "doors-x-via-south"
  },
  "rooms-up-via-relation": {
   "id": "rooms-up-via-relation",
   "name": "rooms-up-via-relation",
   "style": "ManyToOne",
   "src": "rooms-up-via",
   "tgt": "doors-x-via-up"
  },
  "rooms-west-via-relation": {
   "id": "rooms-west-via-relation",
   "name": "rooms-west-via-relation",
   "style": "ManyToOne",
   "src": "rooms-west-via",
   "tgt": "doors-x-via-west"
  },
  "supporters-contents-relation": {
   "id": "supporters-contents-relation",
   "name": "supporters-contents-relation",
   "style": "OneToMany",
   "src": "supporters-contents",
   "tgt": "objects-support"
  }
 },
 "plurals": {
  "actor": "actors",
  "animal": "animals",
  "canister": "canisters",
  "container": "containers",
  "conversation global": "conversation globals",
  "device": "devices",
  "direction": "directions",
  "door": "doors",
  "fact": "facts",
  "following quip": "following quips",
  "kind": "kinds",
  "object": "objects",
  "opener": "openers",
  "pending quip": "pending quips",
  "prop": "props",
  "quip": "quips",
  "quip requirement": "quip requirements",
  "room": "rooms",
  "status bar instance": "status bar instances",
  "story": "stories",
  "supporter": "supporters"
 }
}`
