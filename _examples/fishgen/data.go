package fishgen

var Data = `{
 "Actions": {
  "acquire-it": {
   "Id": "acquire-it",
   "Name": "acquire it",
   "EventId": "acquiring-it",
   "NounTypes": [
    "actors",
    "props"
   ],
   "DefaultActions": [
    "a38fd7888ecbfb8e121d362dbfed"
   ]
  },
  "attack-it": {
   "Id": "attack-it",
   "Name": "attack it",
   "EventId": "attacking-it",
   "NounTypes": [
    "actors",
    "objects"
   ],
   "DefaultActions": [
    "e1067e858bff37ff95ed58744c3655e0"
   ]
  },
  "be-acquired": {
   "Id": "be-acquired",
   "Name": "be acquired",
   "EventId": "being-acquired",
   "NounTypes": [
    "props",
    "actors"
   ],
   "DefaultActions": [
    "aa0f3ad3144d3b009591e22f4f7f9aa2"
   ]
  },
  "be-closed-by": {
   "Id": "be-closed-by",
   "Name": "be closed by",
   "EventId": "being-closed-by",
   "NounTypes": [
    "props",
    "actors"
   ],
   "DefaultActions": [
    "cddd2b508635fcd0fcc15232ddfee0d9"
   ]
  },
  "be-discussed": {
   "Id": "be-discussed",
   "Name": "be discussed",
   "EventId": "being-discussed",
   "NounTypes": [
    "quips",
    "actors"
   ],
   "DefaultActions": [
    "e8ea14871dd41a378b7d216c974bfdd0"
   ]
  },
  "be-examined": {
   "Id": "be-examined",
   "Name": "be examined",
   "EventId": "being-examined",
   "NounTypes": [
    "objects",
    "actors"
   ],
   "DefaultActions": [
    "f5982891e7ebdc570e12599f257d70ea"
   ]
  },
  "be-inserted": {
   "Id": "be-inserted",
   "Name": "be inserted",
   "EventId": "being-inserted",
   "NounTypes": [
    "props",
    "actors",
    "containers"
   ],
   "DefaultActions": [
    "e8b94915f1dace4e2b2f59f3a9162c"
   ]
  },
  "be-opened-by": {
   "Id": "be-opened-by",
   "Name": "be opened by",
   "EventId": "being-opened-by",
   "NounTypes": [
    "props",
    "actors"
   ],
   "DefaultActions": [
    "cce18d15e14eddefbd4c7741dce6f5ce"
   ]
  },
  "be-passed-through": {
   "Id": "be-passed-through",
   "Name": "be passed through",
   "EventId": "being-passed-through",
   "NounTypes": [
    "doors",
    "actors"
   ],
   "DefaultActions": [
    "fe407e25d67b995683a4497a99da6e5e"
   ]
  },
  "close-it": {
   "Id": "close-it",
   "Name": "close it",
   "EventId": "closing-it",
   "NounTypes": [
    "actors",
    "props"
   ],
   "DefaultActions": [
    "a028e313823833dd7dad493203bd86fc"
   ]
  },
  "commence": {
   "Id": "commence",
   "Name": "commence",
   "EventId": "commencing",
   "NounTypes": [
    "stories"
   ],
   "DefaultActions": [
    "be2fb26f4dfcdf76f7d89ab1a25c254"
   ]
  },
  "comment": {
   "Id": "comment",
   "Name": "comment",
   "EventId": "commenting",
   "NounTypes": [
    "actors",
    "quips"
   ],
   "DefaultActions": [
    "efb9c6008baa5f297da22b6b88a69ae"
   ]
  },
  "debug-contents": {
   "Id": "debug-contents",
   "Name": "debug contents",
   "EventId": "debugging-contents",
   "NounTypes": [
    "actors",
    "objects"
   ],
   "DefaultActions": [
    "e27d852c08e661f57fed7a37fb44b6"
   ]
  },
  "debug-direct-parent": {
   "Id": "debug-direct-parent",
   "Name": "debug direct parent",
   "EventId": "debugging-direct-parent",
   "NounTypes": [
    "actors",
    "objects"
   ],
   "DefaultActions": [
    "bb4ae3b5d6669a4eb5d919a52800816"
   ]
  },
  "debug-room-contents": {
   "Id": "debug-room-contents",
   "Name": "debug room contents",
   "EventId": "debugging-room-contents",
   "NounTypes": [
    "actors"
   ],
   "DefaultActions": [
    "efb794467ebba90fb9d54e3f0332632e"
   ]
  },
  "depart": {
   "Id": "depart",
   "Name": "depart",
   "EventId": "departing",
   "NounTypes": [
    "actors"
   ],
   "DefaultActions": [
    "ab0570189d883833d088ab9206526dc3"
   ]
  },
  "describe-the-first-room": {
   "Id": "describe-the-first-room",
   "Name": "describe the first room",
   "EventId": "describing-the-first-room",
   "NounTypes": [
    "stories",
    "rooms"
   ],
   "DefaultActions": [
    "d4983e203d6eeadd7b989e46d6"
   ]
  },
  "discuss": {
   "Id": "discuss",
   "Name": "discuss",
   "EventId": "discussing",
   "NounTypes": [
    "actors",
    "quips"
   ],
   "DefaultActions": [
    "b703927b07424ab1a58cc61fc9287c"
   ]
  },
  "eat-it": {
   "Id": "eat-it",
   "Name": "eat it",
   "EventId": "eating-it",
   "NounTypes": [
    "actors",
    "props"
   ],
   "DefaultActions": [
    "dc5f2a70d4f8b2f822f0d4f669c330b"
   ]
  },
  "end-the-story": {
   "Id": "end-the-story",
   "Name": "end the story",
   "EventId": "ending-the-story",
   "NounTypes": [
    "stories"
   ],
   "DefaultActions": [
    "a1157c67ba6d0bacc2a189d6ac0f88"
   ]
  },
  "end-turn": {
   "Id": "end-turn",
   "Name": "end turn",
   "EventId": "ending-the-turn",
   "NounTypes": [
    "stories"
   ],
   "DefaultActions": [
    "e549441d59279c1801dcd222a57f"
   ]
  },
  "examine-it": {
   "Id": "examine-it",
   "Name": "examine it",
   "EventId": "examining-it",
   "NounTypes": [
    "actors",
    "objects"
   ],
   "DefaultActions": [
    "f3d54a0c2d53c10de008fb6ab276670f"
   ]
  },
  "feed-it": {
   "Id": "feed-it",
   "Name": "feed it",
   "EventId": "feeding-it",
   "NounTypes": [
    "actors",
    "objects"
   ],
   "DefaultActions": [
    "fd94e1c858c8c951a08e2bf2fad194e5"
   ]
  },
  "give-it-to": {
   "Id": "give-it-to",
   "Name": "give it to",
   "EventId": "giving-it-to",
   "NounTypes": [
    "actors",
    "actors",
    "props"
   ],
   "DefaultActions": [
    "c1cd32951100ecf601c99e78a416faa"
   ]
  },
  "go-through-it": {
   "Id": "go-through-it",
   "Name": "go through it",
   "EventId": "going-through-it",
   "NounTypes": [
    "actors",
    "doors"
   ],
   "DefaultActions": [
    "f147fff0a14762b6527a8f71240cdf00"
   ]
  },
  "go-to": {
   "Id": "go-to",
   "Name": "go to",
   "EventId": "going-to",
   "NounTypes": [
    "actors",
    "directions"
   ],
   "DefaultActions": [
    "a279a3ec2a6f0098be46df9a674f96f"
   ]
  },
  "greet": {
   "Id": "greet",
   "Name": "greet",
   "EventId": "greeting",
   "NounTypes": [
    "actors",
    "actors"
   ],
   "DefaultActions": [
    "d73f8083c58056e0b106f0330cabd14e"
   ]
  },
  "impress": {
   "Id": "impress",
   "Name": "impress",
   "EventId": "impressing",
   "NounTypes": [
    "actors"
   ],
   "DefaultActions": [
    "a4e06ed8921a1c0e863719f30eed2"
   ]
  },
  "insert-it-into": {
   "Id": "insert-it-into",
   "Name": "insert it into",
   "EventId": "inserting-it-into",
   "NounTypes": [
    "actors",
    "containers",
    "props"
   ],
   "DefaultActions": [
    "e7456c6895a9da6aa"
   ]
  },
  "jump": {
   "Id": "jump",
   "Name": "jump",
   "EventId": "jumping",
   "NounTypes": [
    "actors"
   ],
   "DefaultActions": [
    "c7e571f4cb23c1766247c3657077f3b9"
   ]
  },
  "kiss-it": {
   "Id": "kiss-it",
   "Name": "kiss it",
   "EventId": "kissing-it",
   "NounTypes": [
    "actors",
    "objects"
   ],
   "DefaultActions": [
    "dbc0e2eebefb53eb032ffd9e2de"
   ]
  },
  "listen": {
   "Id": "listen",
   "Name": "listen",
   "EventId": "listening",
   "NounTypes": [
    "actors"
   ],
   "DefaultActions": [
    "fce05a152a90e9075438ca19ce7c6"
   ]
  },
  "listen-to": {
   "Id": "listen-to",
   "Name": "listen to",
   "EventId": "listening-to",
   "NounTypes": [
    "actors",
    "kinds"
   ],
   "DefaultActions": [
    "fa673c29457c2a9ac47c49658685"
   ]
  },
  "look": {
   "Id": "look",
   "Name": "look",
   "EventId": "looking",
   "NounTypes": [
    "actors"
   ],
   "DefaultActions": [
    "dba636750406cd8d7bd0da3"
   ]
  },
  "look-under-it": {
   "Id": "look-under-it",
   "Name": "look under it",
   "EventId": "looking-under-it",
   "NounTypes": [
    "actors",
    "objects"
   ],
   "DefaultActions": [
    "b0d267fe6bec1f08710b46623fe815b"
   ]
  },
  "open-it": {
   "Id": "open-it",
   "Name": "open it",
   "EventId": "opening-it",
   "NounTypes": [
    "actors",
    "props"
   ],
   "DefaultActions": [
    "b6e604472836e917e8ad1dd6ba2b93"
   ]
  },
  "parse-player-input": {
   "Id": "parse-player-input",
   "Name": "parse player input",
   "EventId": "parsing-player-input",
   "NounTypes": [
    "stories"
   ],
   "DefaultActions": null
  },
  "print-contents": {
   "Id": "print-contents",
   "Name": "print contents",
   "EventId": "printing-contents",
   "NounTypes": [
    "objects"
   ],
   "DefaultActions": null
  },
  "print-conversation-choices": {
   "Id": "print-conversation-choices",
   "Name": "print conversation choices",
   "EventId": "printing-conversation-choices",
   "NounTypes": [
    "actors",
    "actors"
   ],
   "DefaultActions": [
    "e2d83166d1ece71900065119ffe42599"
   ]
  },
  "print-description": {
   "Id": "print-description",
   "Name": "print description",
   "EventId": "describing",
   "NounTypes": [
    "objects"
   ],
   "DefaultActions": [
    "cabe59e7088cb62f7dceba9b137394"
   ]
  },
  "print-details": {
   "Id": "print-details",
   "Name": "print details",
   "EventId": "printing-details",
   "NounTypes": [
    "objects"
   ],
   "DefaultActions": [
    "e5d3985209c83582a8cec2119c5f808"
   ]
  },
  "print-name": {
   "Id": "print-name",
   "Name": "print name",
   "EventId": "printing-name-text",
   "NounTypes": [
    "objects"
   ],
   "DefaultActions": [
    "f3a039df8a36cf28033c549d1d"
   ]
  },
  "print-the-banner": {
   "Id": "print-the-banner",
   "Name": "print the banner",
   "EventId": "printing-the-banner",
   "NounTypes": [
    "stories"
   ],
   "DefaultActions": [
    "f80bca4a7a24cc9ee42067363ecda08"
   ]
  },
  "put-it-onto": {
   "Id": "put-it-onto",
   "Name": "put it onto",
   "EventId": "putting-it-onto",
   "NounTypes": [
    "actors",
    "supporters",
    "props"
   ],
   "DefaultActions": [
    "c7ebd1166c4885183f433a6de51411e1"
   ]
  },
  "receive-insertion": {
   "Id": "receive-insertion",
   "Name": "receive insertion",
   "EventId": "receiving-insertion",
   "NounTypes": [
    "containers",
    "props",
    "actors"
   ],
   "DefaultActions": [
    "bd9ceebcdcc7ebb57be7187da0b3916"
   ]
  },
  "report-already-closed": {
   "Id": "report-already-closed",
   "Name": "report already closed",
   "EventId": "reporting-already-closed",
   "NounTypes": [
    "props",
    "actors"
   ],
   "DefaultActions": [
    "a0cf32d8daa137026c62f8494129"
   ]
  },
  "report-already-off": {
   "Id": "report-already-off",
   "Name": "report already off",
   "EventId": "reporting-already-off",
   "NounTypes": [
    "devices",
    "actors"
   ],
   "DefaultActions": [
    "e54bc28e4a6d6ee8f15cde43d041f41"
   ]
  },
  "report-already-on": {
   "Id": "report-already-on",
   "Name": "report already on",
   "EventId": "reporting-already-on",
   "NounTypes": [
    "devices",
    "actors"
   ],
   "DefaultActions": [
    "b44ad3f16c77942ddebfbcbbdb1381c"
   ]
  },
  "report-already-open": {
   "Id": "report-already-open",
   "Name": "report already open",
   "EventId": "reporting-already-opened",
   "NounTypes": [
    "props",
    "actors"
   ],
   "DefaultActions": [
    "f15ec508566495ef319bda0247c3ba"
   ]
  },
  "report-attack": {
   "Id": "report-attack",
   "Name": "report attack",
   "EventId": "reporting-attack",
   "NounTypes": [
    "objects",
    "actors"
   ],
   "DefaultActions": [
    "e092832c33702a363a16fa5d40421bf"
   ]
  },
  "report-comment": {
   "Id": "report-comment",
   "Name": "report comment",
   "EventId": "reporting-comment",
   "NounTypes": [
    "quips",
    "actors"
   ],
   "DefaultActions": [
    "f1487317dd37aafc7fdc5f81a5296e88"
   ]
  },
  "report-currently-closed": {
   "Id": "report-currently-closed",
   "Name": "report currently closed",
   "EventId": "reporting-currently-closed",
   "NounTypes": [
    "doors",
    "actors"
   ],
   "DefaultActions": [
    "fdb1dc013d4b6490c6586654cf3"
   ]
  },
  "report-eat": {
   "Id": "report-eat",
   "Name": "report eat",
   "EventId": "reporting-eat",
   "NounTypes": [
    "props",
    "actors"
   ],
   "DefaultActions": [
    "f4d2a3f2bc4ff9069889a9299b"
   ]
  },
  "report-gave": {
   "Id": "report-gave",
   "Name": "report gave",
   "EventId": "reporting-gave",
   "NounTypes": [
    "props",
    "actors",
    "actors"
   ],
   "DefaultActions": [
    "f820558419c997630a2044117eacd4b9"
   ]
  },
  "report-give": {
   "Id": "report-give",
   "Name": "report give",
   "EventId": "reporting-give",
   "NounTypes": [
    "actors",
    "props",
    "actors"
   ],
   "DefaultActions": [
    "bb5aa921e7bf8a1a2a243bc496da3d37"
   ]
  },
  "report-inoperable": {
   "Id": "report-inoperable",
   "Name": "report inoperable",
   "EventId": "reporting-inoperable",
   "NounTypes": [
    "props"
   ],
   "DefaultActions": [
    "fdcb93cc0324bd63a08f6415d5f62e3"
   ]
  },
  "report-inventory": {
   "Id": "report-inventory",
   "Name": "report inventory",
   "EventId": "reporting-inventory",
   "NounTypes": [
    "actors"
   ],
   "DefaultActions": [
    "fa0bacce6442c0445e80e1fea52eb2"
   ]
  },
  "report-jump": {
   "Id": "report-jump",
   "Name": "report jump",
   "EventId": "reporting-jump",
   "NounTypes": [
    "kinds",
    "actors"
   ],
   "DefaultActions": [
    "fd97981f521bc27a06faf20b0d5053da"
   ]
  },
  "report-kiss": {
   "Id": "report-kiss",
   "Name": "report kiss",
   "EventId": "reporting-kiss",
   "NounTypes": [
    "objects",
    "actors"
   ],
   "DefaultActions": [
    "a03c941d25050d3a9974d9165bb8f1"
   ]
  },
  "report-listen": {
   "Id": "report-listen",
   "Name": "report listen",
   "EventId": "reporting-listen",
   "NounTypes": [
    "kinds",
    "actors"
   ],
   "DefaultActions": [
    "b0b8e208c4cc475e100e14fbcf2f4ae"
   ]
  },
  "report-locked": {
   "Id": "report-locked",
   "Name": "report locked",
   "EventId": "reporting-locked",
   "NounTypes": [
    "props",
    "actors"
   ],
   "DefaultActions": [
    "b53ef45a3e1238a3aa47a2bdc8d9da4"
   ]
  },
  "report-look-under": {
   "Id": "report-look-under",
   "Name": "report look under",
   "EventId": "reporting-look-under",
   "NounTypes": [
    "objects",
    "actors"
   ],
   "DefaultActions": [
    "eecdc7b5555289906c2b415fb201cb9b"
   ]
  },
  "report-not-closeable": {
   "Id": "report-not-closeable",
   "Name": "report not closeable",
   "EventId": "reporting-not-closeable",
   "NounTypes": [
    "props",
    "actors"
   ],
   "DefaultActions": [
    "d621a08c3de87b0ee4cfe7971a9"
   ]
  },
  "report-now-closed": {
   "Id": "report-now-closed",
   "Name": "report now closed",
   "EventId": "reporting-now-closed",
   "NounTypes": [
    "props",
    "actors"
   ],
   "DefaultActions": [
    "a14a3d5d1cca145f13f50ca1309af311"
   ]
  },
  "report-now-off": {
   "Id": "report-now-off",
   "Name": "report now off",
   "EventId": "reporting-now-off",
   "NounTypes": [
    "devices",
    "actors"
   ],
   "DefaultActions": [
    "bc33e94a52fffabe6b20bcc54955ce3"
   ]
  },
  "report-now-on": {
   "Id": "report-now-on",
   "Name": "report now on",
   "EventId": "reporting-now-on",
   "NounTypes": [
    "devices",
    "actors"
   ],
   "DefaultActions": [
    "f13522698fb406b905407ccdfec39"
   ]
  },
  "report-now-open": {
   "Id": "report-now-open",
   "Name": "report now open",
   "EventId": "reporting-now-open",
   "NounTypes": [
    "props",
    "actors"
   ],
   "DefaultActions": [
    "cd40d14902a558771bcc65444bc871f"
   ]
  },
  "report-placed": {
   "Id": "report-placed",
   "Name": "report placed",
   "EventId": "reporting-placed",
   "NounTypes": [
    "props",
    "actors",
    "supporters"
   ],
   "DefaultActions": [
    "a0ca28fd1e1f6003219b799c59e5925"
   ]
  },
  "report-put": {
   "Id": "report-put",
   "Name": "report put",
   "EventId": "reporting-put",
   "NounTypes": [
    "supporters",
    "props",
    "actors"
   ],
   "DefaultActions": [
    "eba487ee4a161c0ed3a629161b612e"
   ]
  },
  "report-search": {
   "Id": "report-search",
   "Name": "report search",
   "EventId": "reporting-search",
   "NounTypes": [
    "props",
    "actors"
   ],
   "DefaultActions": [
    "c346a8302cd8e56c361e834a576fde7"
   ]
  },
  "report-show": {
   "Id": "report-show",
   "Name": "report show",
   "EventId": "reporting-show",
   "NounTypes": [
    "actors",
    "props",
    "actors"
   ],
   "DefaultActions": [
    "ea0713e6af256dfe532c004ff5a23f5"
   ]
  },
  "report-shown": {
   "Id": "report-shown",
   "Name": "report shown",
   "EventId": "reporting-shown",
   "NounTypes": [
    "props",
    "actors",
    "actors"
   ],
   "DefaultActions": [
    "e460183c17f770fec808d132736f24e"
   ]
  },
  "report-smell": {
   "Id": "report-smell",
   "Name": "report smell",
   "EventId": "reporting-smell",
   "NounTypes": [
    "kinds",
    "actors"
   ],
   "DefaultActions": [
    "fe0b565e4a5208b37419400dbbaea13"
   ]
  },
  "report-switch-off": {
   "Id": "report-switch-off",
   "Name": "report switch off",
   "EventId": "reporting-switch-off",
   "NounTypes": [
    "devices"
   ],
   "DefaultActions": [
    "ef7b600f685b87fcd10d15cb21a3"
   ]
  },
  "report-switched-on": {
   "Id": "report-switched-on",
   "Name": "report switched on",
   "EventId": "reporting-switched-on",
   "NounTypes": [
    "devices",
    "actors"
   ],
   "DefaultActions": [
    "b703969791fe6a4b4b43764fe9215df0"
   ]
  },
  "report-take": {
   "Id": "report-take",
   "Name": "report take",
   "EventId": "reporting-take",
   "NounTypes": [
    "props",
    "actors"
   ],
   "DefaultActions": [
    "c9d7e290147df7748f2222edd327fcb2"
   ]
  },
  "report-the-view": {
   "Id": "report-the-view",
   "Name": "report the view",
   "EventId": "reporting-the-view",
   "NounTypes": [
    "rooms"
   ],
   "DefaultActions": [
    "dbd65ae410d57513e4ddb3770accff"
   ]
  },
  "report-unopenable": {
   "Id": "report-unopenable",
   "Name": "report unopenable",
   "EventId": "reporting-unopenable",
   "NounTypes": [
    "props",
    "actors"
   ],
   "DefaultActions": [
    "c76bb26a42aa9719a9a0313b2ec1193d"
   ]
  },
  "report-wear": {
   "Id": "report-wear",
   "Name": "report wear",
   "EventId": "reporting-wear",
   "NounTypes": [
    "props",
    "actors"
   ],
   "DefaultActions": [
    "dba3f28c029c8ebf74fe7ac2d94e0ac0"
   ]
  },
  "search-it": {
   "Id": "search-it",
   "Name": "search it",
   "EventId": "searching-it",
   "NounTypes": [
    "actors",
    "props"
   ],
   "DefaultActions": [
    "ed97f386ce027a2ff676ef7a8a302e1"
   ]
  },
  "set-initial-position": {
   "Id": "set-initial-position",
   "Name": "set initial position",
   "EventId": "setting-initial-position",
   "NounTypes": [
    "stories",
    "actors",
    "rooms"
   ],
   "DefaultActions": [
    "f6b153d8dfbb91bec15a99f842d2"
   ]
  },
  "show-it-to": {
   "Id": "show-it-to",
   "Name": "show it to",
   "EventId": "showing-it-to",
   "NounTypes": [
    "actors",
    "actors",
    "props"
   ],
   "DefaultActions": [
    "b5bd6ab4d4df9616f49967630082b1e"
   ]
  },
  "smell": {
   "Id": "smell",
   "Name": "smell",
   "EventId": "smelling",
   "NounTypes": [
    "actors"
   ],
   "DefaultActions": [
    "bd1fabe270be6ab6a4724cc8f142da1"
   ]
  },
  "smell-it": {
   "Id": "smell-it",
   "Name": "smell it",
   "EventId": "smelling-it",
   "NounTypes": [
    "actors",
    "kinds"
   ],
   "DefaultActions": [
    "f1eb81a12f8e22e8c9964742f7a703"
   ]
  },
  "switch-it-off": {
   "Id": "switch-it-off",
   "Name": "switch it off",
   "EventId": "switching-it-off",
   "NounTypes": [
    "actors",
    "props"
   ],
   "DefaultActions": [
    "d61637c94b5914263b9b2e888c6ee3f0"
   ]
  },
  "switch-it-on": {
   "Id": "switch-it-on",
   "Name": "switch it on",
   "EventId": "switching-it-on",
   "NounTypes": [
    "actors",
    "props"
   ],
   "DefaultActions": [
    "c8f1338cac6671fc5d6920c2c0f7b893"
   ]
  },
  "take-it": {
   "Id": "take-it",
   "Name": "take it",
   "EventId": "taking-it",
   "NounTypes": [
    "actors",
    "props"
   ],
   "DefaultActions": [
    "ad645ba09d5207de6d3ab358466b7355"
   ]
  },
  "wear-it": {
   "Id": "wear-it",
   "Name": "wear it",
   "EventId": "wearing-it",
   "NounTypes": [
    "actors",
    "props"
   ],
   "DefaultActions": [
    "a2afcc3585516a895061346c3d5fa3f"
   ]
  }
 },
 "Classes": {
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
    },
    {
     "id": "actors-not-inputing-dialog",
     "type": "EnumProperty",
     "name": "not inputing dialog"
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
     "id": "containers-contents",
     "type": "PointerProperty",
     "name": "contents",
     "relates": "objects",
     "relation": "containers-contents-relation",
     "many": true
    },
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
    },
    {
     "id": "conversation-globals-queue",
     "type": "PointerProperty",
     "name": "queue",
     "relates": "quips",
     "many": true
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
     "id": "doors-x-via-up",
     "type": "PointerProperty",
     "name": "x-via-up",
     "relates": "rooms",
     "relation": "rooms-up-via-relation",
     "many": true
    },
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
     "id": "objects-whereabouts",
     "type": "PointerProperty",
     "name": "whereabouts",
     "relates": "rooms",
     "relation": "rooms-contents-relation"
    },
    {
     "id": "objects-description",
     "type": "TextProperty",
     "name": "description"
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
     "id": "objects-wearer",
     "type": "PointerProperty",
     "name": "wearer",
     "relates": "actors",
     "relation": "actors-clothing-relation"
    },
    {
     "id": "objects-enclosure",
     "type": "PointerProperty",
     "name": "enclosure",
     "relates": "containers",
     "relation": "containers-contents-relation"
    },
    {
     "id": "objects-unhandled",
     "type": "EnumProperty",
     "name": "unhandled"
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
    },
    {
     "id": "openers-unlocked",
     "type": "EnumProperty",
     "name": "unlocked"
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
     "id": "props-portable",
     "type": "EnumProperty",
     "name": "portable"
    },
    {
     "id": "props-wearable",
     "type": "EnumProperty",
     "name": "wearable"
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
     "id": "quip-requirements-permitted",
     "type": "EnumProperty",
     "name": "permitted"
    },
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
     "id": "quips-restrictive",
     "type": "EnumProperty",
     "name": "restrictive"
    },
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
     "id": "rooms-description",
     "type": "TextProperty",
     "name": "description"
    },
    {
     "id": "rooms-south-via",
     "type": "PointerProperty",
     "name": "south-via",
     "relates": "doors",
     "relation": "rooms-south-via-relation"
    },
    {
     "id": "rooms-up-via",
     "type": "PointerProperty",
     "name": "up-via",
     "relates": "doors",
     "relation": "rooms-up-via-relation"
    },
    {
     "id": "rooms-down-via",
     "type": "PointerProperty",
     "name": "down-via",
     "relates": "doors",
     "relation": "rooms-down-via-relation"
    },
    {
     "id": "rooms-visited",
     "type": "EnumProperty",
     "name": "visited"
    },
    {
     "id": "rooms-east-via",
     "type": "PointerProperty",
     "name": "east-via",
     "relates": "doors",
     "relation": "rooms-east-via-relation"
    },
    {
     "id": "rooms-west-via",
     "type": "PointerProperty",
     "name": "west-via",
     "relates": "doors",
     "relation": "rooms-west-via-relation"
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
    },
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
 "Enumerations": {
  "actors-not-inputing-dialog": {
   "Choices": [
    "not-inputing-dialog",
    "inputing-dialog"
   ]
  },
  "canisters-hidden": {
   "Choices": [
    "hidden",
    "found"
   ]
  },
  "containers-lockable": {
   "Choices": [
    "not-lockable",
    "lockable"
   ]
  },
  "containers-locked": {
   "Choices": [
    "unlocked",
    "locked"
   ]
  },
  "containers-opaque": {
   "Choices": [
    "opaque",
    "transparent"
   ]
  },
  "devices-operable": {
   "Choices": [
    "operable",
    "inoperable"
   ]
  },
  "devices-switched-off": {
   "Choices": [
    "switched-off",
    "switched-on"
   ]
  },
  "following-quips-indirectly-following": {
   "Choices": [
    "indirectly-following",
    "directly-following"
   ]
  },
  "kinds-common-named": {
   "Choices": [
    "common-named",
    "proper-named"
   ]
  },
  "kinds-recollected": {
   "Choices": [
    "not-recollected",
    "recollected"
   ]
  },
  "kinds-singular-named": {
   "Choices": [
    "singular-named",
    "plural-named"
   ]
  },
  "objects-scenery": {
   "Choices": [
    "not-scenery",
    "scenery"
   ]
  },
  "objects-unhandled": {
   "Choices": [
    "unhandled",
    "handled"
   ]
  },
  "openers-hinged": {
   "Choices": [
    "hinged",
    "hingeless"
   ]
  },
  "openers-locakable": {
   "Choices": [
    "not-lockable",
    "locakable"
   ]
  },
  "openers-open": {
   "Choices": [
    "open",
    "closed"
   ]
  },
  "openers-unlocked": {
   "Choices": [
    "unlocked",
    "locked"
   ]
  },
  "pending-quips-immediate": {
   "Choices": [
    "immediate",
    "postponed"
   ]
  },
  "pending-quips-obligatory": {
   "Choices": [
    "obligatory",
    "optional"
   ]
  },
  "props-portable": {
   "Choices": [
    "portable",
    "fixed-in-place"
   ]
  },
  "props-wearable": {
   "Choices": [
    "not-wearable",
    "wearable"
   ]
  },
  "quip-requirements-permitted": {
   "Choices": [
    "permitted",
    "prohibited"
   ]
  },
  "quips-planned": {
   "Choices": [
    "planned",
    "casual"
   ]
  },
  "quips-repeatable": {
   "Choices": [
    "one-time",
    "repeatable"
   ]
  },
  "quips-restrictive": {
   "Choices": [
    "unrestricted",
    "restrictive"
   ]
  },
  "rooms-visited": {
   "Choices": [
    "unvisited",
    "visited"
   ]
  },
  "stories-playing": {
   "Choices": [
    "playing",
    "completed"
   ]
  },
  "stories-scored": {
   "Choices": [
    "unscored",
    "scored"
   ]
  }
 },
 "Events": {
  "acquiring-it": {
   "Id": "acquiring-it",
   "Name": "acquiring it",
   "Capture": null,
   "Bubble": null
  },
  "attacking-it": {
   "Id": "attacking-it",
   "Name": "attacking it",
   "Capture": null,
   "Bubble": null
  },
  "being-acquired": {
   "Id": "being-acquired",
   "Name": "being acquired",
   "Capture": null,
   "Bubble": null
  },
  "being-closed-by": {
   "Id": "being-closed-by",
   "Name": "being closed by",
   "Capture": null,
   "Bubble": null
  },
  "being-discussed": {
   "Id": "being-discussed",
   "Name": "being discussed",
   "Capture": null,
   "Bubble": null
  },
  "being-examined": {
   "Id": "being-examined",
   "Name": "being examined",
   "Capture": [
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
     "Callback": "c67b0b4d7d698d094c906154",
     "Options": 1
    },
    {
     "Instance": "painting",
     "Class": "props",
     "Callback": "ec6f223e7791279fda070c1800163efb",
     "Options": 3
    },
    {
     "Instance": "table",
     "Class": "supporters",
     "Callback": "f1004c0b7326249efe3f65557fabb9f4",
     "Options": 3
    },
    {
     "Instance": "telegram",
     "Class": "props",
     "Callback": "e69b7d9edaeed0a36460207695",
     "Options": 3
    },
    {
     "Instance": "gravel",
     "Class": "props",
     "Callback": "bf01d4248f78ebe97c2271f32858372c",
     "Options": 3
    },
    {
     "Instance": "seaweed",
     "Class": "props",
     "Callback": "dcd93e8528fa26a2d76ff388e12b05",
     "Options": 3
    },
    {
     "Instance": "evil-fish",
     "Class": "animals",
     "Callback": "b5220b2fe33551f0f16fcb86e43b",
     "Options": 3
    },
    {
     "Instance": "bouquet",
     "Class": "props",
     "Callback": "a1f234db8a0aac739ec97f31970dc",
     "Options": 3
    },
    {
     "Instance": "lingerie-bag",
     "Class": "props",
     "Callback": "c490a9a325bf86d652e6f8ae46f",
     "Options": 3
    }
   ],
   "Bubble": [
    {
     "Instance": "aquarium",
     "Class": "containers",
     "Callback": "ce5b032afbe64e94fbb24b1cb91",
     "Options": 1
    }
   ]
  },
  "being-inserted": {
   "Id": "being-inserted",
   "Name": "being inserted",
   "Capture": null,
   "Bubble": [
    {
     "Instance": "fish-food",
     "Class": "canisters",
     "Callback": "f3b6951a9a4afccf925531d9527d552e",
     "Options": 1
    }
   ]
  },
  "being-opened-by": {
   "Id": "being-opened-by",
   "Name": "being opened by",
   "Capture": null,
   "Bubble": [
    {
     "Instance": "fish-food",
     "Class": "canisters",
     "Callback": "e06a20281952b2e90d79adf1a13caa4",
     "Options": 1
    }
   ]
  },
  "being-passed-through": {
   "Id": "being-passed-through",
   "Name": "being passed through",
   "Capture": null,
   "Bubble": null
  },
  "closing-it": {
   "Id": "closing-it",
   "Name": "closing it",
   "Capture": null,
   "Bubble": null
  },
  "commencing": {
   "Id": "commencing",
   "Name": "commencing",
   "Capture": null,
   "Bubble": [
    {
     "Instance": "day-for-fresh-sushi",
     "Class": "stories",
     "Callback": "d973ac60ad949239200a0ded1e608380",
     "Options": 1
    }
   ]
  },
  "commenting": {
   "Id": "commenting",
   "Name": "commenting",
   "Capture": null,
   "Bubble": null
  },
  "debugging-contents": {
   "Id": "debugging-contents",
   "Name": "debugging contents",
   "Capture": null,
   "Bubble": null
  },
  "debugging-direct-parent": {
   "Id": "debugging-direct-parent",
   "Name": "debugging direct parent",
   "Capture": null,
   "Bubble": null
  },
  "debugging-room-contents": {
   "Id": "debugging-room-contents",
   "Name": "debugging room contents",
   "Capture": null,
   "Bubble": null
  },
  "departing": {
   "Id": "departing",
   "Name": "departing",
   "Capture": null,
   "Bubble": null
  },
  "describing": {
   "Id": "describing",
   "Name": "describing",
   "Capture": null,
   "Bubble": null
  },
  "describing-the-first-room": {
   "Id": "describing-the-first-room",
   "Name": "describing the first room",
   "Capture": null,
   "Bubble": null
  },
  "discussing": {
   "Id": "discussing",
   "Name": "discussing",
   "Capture": null,
   "Bubble": null
  },
  "eating-it": {
   "Id": "eating-it",
   "Name": "eating it",
   "Capture": null,
   "Bubble": null
  },
  "ending-the-story": {
   "Id": "ending-the-story",
   "Name": "ending the story",
   "Capture": null,
   "Bubble": null
  },
  "ending-the-turn": {
   "Id": "ending-the-turn",
   "Name": "ending the turn",
   "Capture": [
    {
     "Instance": "",
     "Class": "stories",
     "Callback": "e994180f6a5d75a96cbdd199f4385b7",
     "Options": 1
    }
   ],
   "Bubble": [
    {
     "Instance": "",
     "Class": "stories",
     "Callback": "c937462c89e5bdcf0dd8d0922d98",
     "Options": 1
    },
    {
     "Instance": "",
     "Class": "stories",
     "Callback": "b1e679677e49de176b715c792ab9ef",
     "Options": 1
    }
   ]
  },
  "examining-it": {
   "Id": "examining-it",
   "Name": "examining it",
   "Capture": null,
   "Bubble": null
  },
  "feeding-it": {
   "Id": "feeding-it",
   "Name": "feeding it",
   "Capture": [
    {
     "Instance": "",
     "Class": "actors",
     "Callback": "b3ce757fe1ce243e795a95b8a480",
     "Options": 1
    }
   ],
   "Bubble": [
    {
     "Instance": "",
     "Class": "actors",
     "Callback": "f62ee471906ed55e81c1e249d3a4651",
     "Options": 1
    },
    {
     "Instance": "",
     "Class": "actors",
     "Callback": "f933c558efc4b3edc240d80b1f4f91d3",
     "Options": 1
    }
   ]
  },
  "giving-it-to": {
   "Id": "giving-it-to",
   "Name": "giving it to",
   "Capture": [
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
   ],
   "Bubble": null
  },
  "going-through-it": {
   "Id": "going-through-it",
   "Name": "going through it",
   "Capture": null,
   "Bubble": null
  },
  "going-to": {
   "Id": "going-to",
   "Name": "going to",
   "Capture": null,
   "Bubble": null
  },
  "greeting": {
   "Id": "greeting",
   "Name": "greeting",
   "Capture": null,
   "Bubble": null
  },
  "impressing": {
   "Id": "impressing",
   "Name": "impressing",
   "Capture": null,
   "Bubble": null
  },
  "inserting-it-into": {
   "Id": "inserting-it-into",
   "Name": "inserting it into",
   "Capture": [
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
   ],
   "Bubble": null
  },
  "jumping": {
   "Id": "jumping",
   "Name": "jumping",
   "Capture": null,
   "Bubble": [
    {
     "Instance": "player",
     "Class": "actors",
     "Callback": "d3aad960bedd830691389807b6897",
     "Options": 1
    }
   ]
  },
  "kissing-it": {
   "Id": "kissing-it",
   "Name": "kissing it",
   "Capture": [
    {
     "Instance": "",
     "Class": "actors",
     "Callback": "f9baafd0c928bb023c1ae00466c0464f",
     "Options": 0
    }
   ],
   "Bubble": null
  },
  "listening": {
   "Id": "listening",
   "Name": "listening",
   "Capture": null,
   "Bubble": null
  },
  "listening-to": {
   "Id": "listening-to",
   "Name": "listening to",
   "Capture": null,
   "Bubble": null
  },
  "looking": {
   "Id": "looking",
   "Name": "looking",
   "Capture": null,
   "Bubble": null
  },
  "looking-under-it": {
   "Id": "looking-under-it",
   "Name": "looking under it",
   "Capture": null,
   "Bubble": null
  },
  "opening-it": {
   "Id": "opening-it",
   "Name": "opening it",
   "Capture": null,
   "Bubble": null
  },
  "parsing-player-input": {
   "Id": "parsing-player-input",
   "Name": "parsing player input",
   "Capture": null,
   "Bubble": [
    {
     "Instance": "",
     "Class": "stories",
     "Callback": "c5f3dd965e3297200df816e37587df",
     "Options": 1
    }
   ]
  },
  "printing-contents": {
   "Id": "printing-contents",
   "Name": "printing contents",
   "Capture": null,
   "Bubble": [
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
   "Id": "printing-conversation-choices",
   "Name": "printing conversation choices",
   "Capture": null,
   "Bubble": null
  },
  "printing-details": {
   "Id": "printing-details",
   "Name": "printing details",
   "Capture": null,
   "Bubble": [
    {
     "Instance": "studio",
     "Class": "rooms",
     "Callback": "d95e07eb79eaf3fa2f3302d5ec0bf09",
     "Options": 1
    },
    {
     "Instance": "window",
     "Class": "openers",
     "Callback": "bf33e6851af7c5602ee19fc40c4ac",
     "Options": 1
    }
   ]
  },
  "printing-name-text": {
   "Id": "printing-name-text",
   "Name": "printing name text",
   "Capture": null,
   "Bubble": [
    {
     "Instance": "",
     "Class": "containers",
     "Callback": "ec75951720bb443a92e928bf812c3",
     "Options": 1
    },
    {
     "Instance": "",
     "Class": "doors",
     "Callback": "a23027802673435b37f0098b3dcb5cfc",
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
   "Id": "printing-the-banner",
   "Name": "printing the banner",
   "Capture": null,
   "Bubble": null
  },
  "putting-it-onto": {
   "Id": "putting-it-onto",
   "Name": "putting it onto",
   "Capture": [
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
   ],
   "Bubble": null
  },
  "receiving-insertion": {
   "Id": "receiving-insertion",
   "Name": "receiving insertion",
   "Capture": [
    {
     "Instance": "",
     "Class": "containers",
     "Callback": "ff8957b596f3ee3ccd61c07445b3e2c8",
     "Options": 0
    }
   ],
   "Bubble": [
    {
     "Instance": "vase",
     "Class": "containers",
     "Callback": "ba81c1a9607db24dfe523984dce6a",
     "Options": 1
    }
   ]
  },
  "reporting-already-closed": {
   "Id": "reporting-already-closed",
   "Name": "reporting already closed",
   "Capture": null,
   "Bubble": null
  },
  "reporting-already-off": {
   "Id": "reporting-already-off",
   "Name": "reporting already off",
   "Capture": null,
   "Bubble": null
  },
  "reporting-already-on": {
   "Id": "reporting-already-on",
   "Name": "reporting already on",
   "Capture": null,
   "Bubble": null
  },
  "reporting-already-opened": {
   "Id": "reporting-already-opened",
   "Name": "reporting already opened",
   "Capture": null,
   "Bubble": null
  },
  "reporting-attack": {
   "Id": "reporting-attack",
   "Name": "reporting attack",
   "Capture": null,
   "Bubble": [
    {
     "Instance": "evil-fish",
     "Class": "animals",
     "Callback": "bfe132582370c1a1242feeb43c5208",
     "Options": 1
    }
   ]
  },
  "reporting-comment": {
   "Id": "reporting-comment",
   "Name": "reporting comment",
   "Capture": null,
   "Bubble": null
  },
  "reporting-currently-closed": {
   "Id": "reporting-currently-closed",
   "Name": "reporting currently closed",
   "Capture": null,
   "Bubble": null
  },
  "reporting-eat": {
   "Id": "reporting-eat",
   "Name": "reporting eat",
   "Capture": null,
   "Bubble": null
  },
  "reporting-gave": {
   "Id": "reporting-gave",
   "Name": "reporting gave",
   "Capture": null,
   "Bubble": [
    {
     "Instance": "fish-food",
     "Class": "canisters",
     "Callback": "e00bf959e49b44bce2d5e2ee91dd99e",
     "Options": 1
    }
   ]
  },
  "reporting-give": {
   "Id": "reporting-give",
   "Name": "reporting give",
   "Capture": null,
   "Bubble": null
  },
  "reporting-inoperable": {
   "Id": "reporting-inoperable",
   "Name": "reporting inoperable",
   "Capture": null,
   "Bubble": null
  },
  "reporting-inventory": {
   "Id": "reporting-inventory",
   "Name": "reporting inventory",
   "Capture": null,
   "Bubble": null
  },
  "reporting-jump": {
   "Id": "reporting-jump",
   "Name": "reporting jump",
   "Capture": null,
   "Bubble": null
  },
  "reporting-kiss": {
   "Id": "reporting-kiss",
   "Name": "reporting kiss",
   "Capture": null,
   "Bubble": [
    {
     "Instance": "evil-fish",
     "Class": "animals",
     "Callback": "dc02c3a20d531a66f5db42562bb2cc2",
     "Options": 1
    }
   ]
  },
  "reporting-listen": {
   "Id": "reporting-listen",
   "Name": "reporting listen",
   "Capture": null,
   "Bubble": null
  },
  "reporting-locked": {
   "Id": "reporting-locked",
   "Name": "reporting locked",
   "Capture": null,
   "Bubble": null
  },
  "reporting-look-under": {
   "Id": "reporting-look-under",
   "Name": "reporting look under",
   "Capture": null,
   "Bubble": [
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
     "Callback": "c6a67fbc6f638cb3d52658d053f1",
     "Options": 1
    }
   ]
  },
  "reporting-not-closeable": {
   "Id": "reporting-not-closeable",
   "Name": "reporting not closeable",
   "Capture": null,
   "Bubble": null
  },
  "reporting-now-closed": {
   "Id": "reporting-now-closed",
   "Name": "reporting now closed",
   "Capture": null,
   "Bubble": [
    {
     "Instance": "cabinet",
     "Class": "containers",
     "Callback": "df1fbacb62e40599f5b43f7e52eb",
     "Options": 1
    }
   ]
  },
  "reporting-now-off": {
   "Id": "reporting-now-off",
   "Name": "reporting now off",
   "Capture": null,
   "Bubble": null
  },
  "reporting-now-on": {
   "Id": "reporting-now-on",
   "Name": "reporting now on",
   "Capture": null,
   "Bubble": null
  },
  "reporting-now-open": {
   "Id": "reporting-now-open",
   "Name": "reporting now open",
   "Capture": null,
   "Bubble": [
    {
     "Instance": "cabinet",
     "Class": "containers",
     "Callback": "c239f2d147803170152e8566a7",
     "Options": 1
    },
    {
     "Instance": "window",
     "Class": "openers",
     "Callback": "dbb2267e60b061fa54065b480b4e3e3",
     "Options": 1
    }
   ]
  },
  "reporting-placed": {
   "Id": "reporting-placed",
   "Name": "reporting placed",
   "Capture": null,
   "Bubble": null
  },
  "reporting-put": {
   "Id": "reporting-put",
   "Name": "reporting put",
   "Capture": null,
   "Bubble": null
  },
  "reporting-search": {
   "Id": "reporting-search",
   "Name": "reporting search",
   "Capture": null,
   "Bubble": [
    {
     "Instance": "cloths",
     "Class": "props",
     "Callback": "c5351efbb7da045f6635e48a3e1958d",
     "Options": 1
    }
   ]
  },
  "reporting-show": {
   "Id": "reporting-show",
   "Name": "reporting show",
   "Capture": null,
   "Bubble": null
  },
  "reporting-shown": {
   "Id": "reporting-shown",
   "Name": "reporting shown",
   "Capture": null,
   "Bubble": [
    {
     "Instance": "cloths",
     "Class": "props",
     "Callback": "a795e203a59b3f39adb4b28ca8eef4c3",
     "Options": 1
    }
   ]
  },
  "reporting-smell": {
   "Id": "reporting-smell",
   "Name": "reporting smell",
   "Capture": null,
   "Bubble": [
    {
     "Instance": "studio",
     "Class": "rooms",
     "Callback": "a98fd9211e9759d489fb4522c66d852f",
     "Options": 1
    },
    {
     "Instance": "bouquet",
     "Class": "props",
     "Callback": "b2256ccb5b70980396f772a35a7571ec",
     "Options": 1
    }
   ]
  },
  "reporting-switch-off": {
   "Id": "reporting-switch-off",
   "Name": "reporting switch off",
   "Capture": null,
   "Bubble": null
  },
  "reporting-switched-on": {
   "Id": "reporting-switched-on",
   "Name": "reporting switched on",
   "Capture": null,
   "Bubble": null
  },
  "reporting-take": {
   "Id": "reporting-take",
   "Name": "reporting take",
   "Capture": [
    {
     "Instance": "",
     "Class": "doors",
     "Callback": "d6fb6189d5ad739d1f229b238d88c",
     "Options": 1
    }
   ],
   "Bubble": [
    {
     "Instance": "paints",
     "Class": "props",
     "Callback": "d87833ece051f1acbf636ab531b13708",
     "Options": 1
    },
    {
     "Instance": "painting",
     "Class": "props",
     "Callback": "b6e04bb9c5925c8ab1d90d67b4ca5",
     "Options": 1
    },
    {
     "Instance": "evil-fish",
     "Class": "animals",
     "Callback": "d74b3b09863c03d984b9a52db6b38447",
     "Options": 1
    }
   ]
  },
  "reporting-the-view": {
   "Id": "reporting-the-view",
   "Name": "reporting the view",
   "Capture": null,
   "Bubble": null
  },
  "reporting-unopenable": {
   "Id": "reporting-unopenable",
   "Name": "reporting unopenable",
   "Capture": null,
   "Bubble": null
  },
  "reporting-wear": {
   "Id": "reporting-wear",
   "Name": "reporting wear",
   "Capture": null,
   "Bubble": null
  },
  "searching-it": {
   "Id": "searching-it",
   "Name": "searching it",
   "Capture": null,
   "Bubble": null
  },
  "setting-initial-position": {
   "Id": "setting-initial-position",
   "Name": "setting initial position",
   "Capture": null,
   "Bubble": null
  },
  "showing-it-to": {
   "Id": "showing-it-to",
   "Name": "showing it to",
   "Capture": [
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
   ],
   "Bubble": null
  },
  "smelling": {
   "Id": "smelling",
   "Name": "smelling",
   "Capture": null,
   "Bubble": null
  },
  "smelling-it": {
   "Id": "smelling-it",
   "Name": "smelling it",
   "Capture": null,
   "Bubble": null
  },
  "switching-it-off": {
   "Id": "switching-it-off",
   "Name": "switching it off",
   "Capture": null,
   "Bubble": null
  },
  "switching-it-on": {
   "Id": "switching-it-on",
   "Name": "switching it on",
   "Capture": null,
   "Bubble": null
  },
  "taking-it": {
   "Id": "taking-it",
   "Name": "taking it",
   "Capture": null,
   "Bubble": null
  },
  "wearing-it": {
   "Id": "wearing-it",
   "Name": "wearing it",
   "Capture": null,
   "Bubble": null
  }
 },
 "Instances": {
  "aquarium": {
   "id": "aquarium",
   "type": "containers",
   "name": "aquarium",
   "Values": {
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
   "Values": {
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
   "Values": {
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
   "Values": {
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
   "Values": {
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
   "Values": {
    "kinds-name": "closedCabinet"
   }
  },
  "cloths": {
   "id": "cloths",
   "type": "props",
   "name": "cloths",
   "Values": {
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
   "Values": {
    "kinds-name": "conversation"
   }
  },
  "day-for-fresh-sushi": {
   "id": "day-for-fresh-sushi",
   "type": "stories",
   "name": "Day For Fresh Sushi",
   "Values": {
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
   "Values": {
    "directions-opposite": "up",
    "directions-x-opposite": "up",
    "kinds-name": "down"
   }
  },
  "easel": {
   "id": "easel",
   "type": "supporters",
   "name": "easel",
   "Values": {
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
   "Values": {
    "directions-opposite": "west",
    "directions-x-opposite": "west",
    "kinds-name": "east"
   }
  },
  "evil-fish": {
   "id": "evil-fish",
   "type": "animals",
   "name": "evil fish",
   "Values": {
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
   "Values": {
    "kinds-name": "examinedBagOnce"
   }
  },
  "examined-bag-twice": {
   "id": "examined-bag-twice",
   "type": "facts",
   "name": "examinedBagTwice",
   "Values": {
    "kinds-name": "examinedBagTwice"
   }
  },
  "examined-bouquet": {
   "id": "examined-bouquet",
   "type": "facts",
   "name": "examinedBouquet",
   "Values": {
    "kinds-name": "examinedBouquet"
   }
  },
  "examined-cloths": {
   "id": "examined-cloths",
   "type": "facts",
   "name": "examinedCloths",
   "Values": {
    "kinds-name": "examinedCloths"
   }
  },
  "examined-fish-once": {
   "id": "examined-fish-once",
   "type": "facts",
   "name": "examinedFishOnce",
   "Values": {
    "kinds-name": "examinedFishOnce"
   }
  },
  "examined-fish-twice": {
   "id": "examined-fish-twice",
   "type": "facts",
   "name": "examinedFishTwice",
   "Values": {
    "kinds-name": "examinedFishTwice"
   }
  },
  "examined-gravel": {
   "id": "examined-gravel",
   "type": "facts",
   "name": "examinedGravel",
   "Values": {
    "kinds-name": "examinedGravel"
   }
  },
  "examined-painting": {
   "id": "examined-painting",
   "type": "facts",
   "name": "examinedPainting",
   "Values": {
    "kinds-name": "examinedPainting"
   }
  },
  "examined-paints": {
   "id": "examined-paints",
   "type": "facts",
   "name": "examinedPaints",
   "Values": {
    "kinds-name": "examinedPaints"
   }
  },
  "examined-seaweed": {
   "id": "examined-seaweed",
   "type": "facts",
   "name": "examinedSeaweed",
   "Values": {
    "kinds-name": "examinedSeaweed"
   }
  },
  "examined-telegraph": {
   "id": "examined-telegraph",
   "type": "facts",
   "name": "examinedTelegraph",
   "Values": {
    "kinds-name": "examinedTelegraph"
   }
  },
  "fish-food": {
   "id": "fish-food",
   "type": "canisters",
   "name": "fish food",
   "Values": {
    "canisters-hidden": "hidden",
    "containers-opaque": "opaque",
    "kinds-name": "fish food",
    "kinds-singular-named": "plural-named",
    "objects-description": "A vehemently orange canister of fish food.",
    "objects-enclosure": "",
    "objects-owner": "",
    "objects-support": "",
    "objects-wearer": "",
    "objects-whereabouts": ""
   }
  },
  "gravel": {
   "id": "gravel",
   "type": "props",
   "name": "gravel",
   "Values": {
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
   "Values": {
    "kinds-name": "insertedFlowers"
   }
  },
  "lingerie-bag": {
   "id": "lingerie-bag",
   "type": "props",
   "name": "lingerie bag",
   "Values": {
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
   "Values": {
    "kinds-name": "lookedUnderCabinet"
   }
  },
  "looked-under-table": {
   "id": "looked-under-table",
   "type": "facts",
   "name": "lookedUnderTable",
   "Values": {
    "kinds-name": "lookedUnderTable"
   }
  },
  "north": {
   "id": "north",
   "type": "directions",
   "name": "north",
   "Values": {
    "directions-opposite": "south",
    "directions-x-opposite": "south",
    "kinds-name": "north"
   }
  },
  "opened-cabinet": {
   "id": "opened-cabinet",
   "type": "facts",
   "name": "openedCabinet",
   "Values": {
    "kinds-name": "openedCabinet"
   }
  },
  "opened-window": {
   "id": "opened-window",
   "type": "facts",
   "name": "openedWindow",
   "Values": {
    "kinds-name": "openedWindow"
   }
  },
  "painting": {
   "id": "painting",
   "type": "props",
   "name": "painting",
   "Values": {
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
   "Values": {
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
   "Values": {
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
   "Values": {
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
   "Values": {
    "kinds-name": "smelledBouquet"
   }
  },
  "south": {
   "id": "south",
   "type": "directions",
   "name": "south",
   "Values": {
    "directions-opposite": "north",
    "directions-x-opposite": "north",
    "kinds-name": "south"
   }
  },
  "status-bar": {
   "id": "status-bar",
   "type": "status-bar-instances",
   "name": "status bar",
   "Values": {
    "kinds-name": "status bar"
   }
  },
  "studio": {
   "id": "studio",
   "type": "rooms",
   "name": "studio",
   "Values": {
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
   "Values": {
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
   "Values": {
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
   "Values": {
    "kinds-name": "tookPaints"
   }
  },
  "up": {
   "id": "up",
   "type": "directions",
   "name": "up",
   "Values": {
    "directions-opposite": "down",
    "directions-x-opposite": "down",
    "kinds-name": "up"
   }
  },
  "vase": {
   "id": "vase",
   "type": "containers",
   "name": "vase",
   "Values": {
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
   "Values": {
    "directions-opposite": "east",
    "directions-x-opposite": "east",
    "kinds-name": "west"
   }
  },
  "window": {
   "id": "window",
   "type": "openers",
   "name": "window",
   "Values": {
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
 "Aliases": {
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
   "examined-bouquet",
   "smelled-bouquet",
   "bouquet"
  ],
  "britney": [
   "britney",
   "britney"
  ],
  "cabinet": [
   "cabinet",
   "looked-under-cabinet",
   "closed-cabinet",
   "cabinet",
   "opened-cabinet"
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
   "cloths",
   "examined-cloths"
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
   "examined-bouquet",
   "examined-fish-twice",
   "examined-bag-twice",
   "examined-telegraph",
   "examined-fish-once",
   "examined-bag-once",
   "examined-painting",
   "examined-seaweed",
   "examined-paints",
   "examined-gravel",
   "examined-cloths"
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
   "examined-fish-twice",
   "examined-fish-once",
   "fish-food"
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
   "examined-fish-once",
   "examined-bag-once"
  ],
  "opened": [
   "opened-cabinet",
   "opened-window"
  ],
  "opened cabinet": [
   "opened-cabinet"
  ],
  "opened window": [
   "opened-window"
  ],
  "painting": [
   "painting",
   "painting",
   "examined-painting"
  ],
  "paints": [
   "paints",
   "paints",
   "took-paints",
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
   "examined-seaweed",
   "seaweed"
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
   "examined-fish-twice",
   "examined-bag-twice"
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
 "ParserActions": [
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
 "Relations": {
  "actors-clothing-relation": {
   "Id": "actors-clothing-relation",
   "Name": "actors-clothing-relation",
   "Style": "OneToMany",
   "Source": "actors-clothing",
   "Target": "objects-wearer"
  },
  "actors-inventory-relation": {
   "Id": "actors-inventory-relation",
   "Name": "actors-inventory-relation",
   "Style": "OneToMany",
   "Source": "actors-inventory",
   "Target": "objects-owner"
  },
  "containers-contents-relation": {
   "Id": "containers-contents-relation",
   "Name": "containers-contents-relation",
   "Style": "OneToMany",
   "Source": "containers-contents",
   "Target": "objects-enclosure"
  },
  "directions-opposite-relation": {
   "Id": "directions-opposite-relation",
   "Name": "directions-opposite-relation",
   "Style": "OneToOne",
   "Source": "directions-opposite",
   "Target": "directions-x-opposite"
  },
  "doors-destination-relation": {
   "Id": "doors-destination-relation",
   "Name": "doors-destination-relation",
   "Style": "ManyToOne",
   "Source": "doors-destination",
   "Target": "doors-sources"
  },
  "rooms-contents-relation": {
   "Id": "rooms-contents-relation",
   "Name": "rooms-contents-relation",
   "Style": "OneToMany",
   "Source": "rooms-contents",
   "Target": "objects-whereabouts"
  },
  "rooms-down-via-relation": {
   "Id": "rooms-down-via-relation",
   "Name": "rooms-down-via-relation",
   "Style": "ManyToOne",
   "Source": "rooms-down-via",
   "Target": "doors-x-via-down"
  },
  "rooms-east-via-relation": {
   "Id": "rooms-east-via-relation",
   "Name": "rooms-east-via-relation",
   "Style": "ManyToOne",
   "Source": "rooms-east-via",
   "Target": "doors-x-via-east"
  },
  "rooms-north-via-relation": {
   "Id": "rooms-north-via-relation",
   "Name": "rooms-north-via-relation",
   "Style": "ManyToOne",
   "Source": "rooms-north-via",
   "Target": "doors-x-via-north"
  },
  "rooms-south-via-relation": {
   "Id": "rooms-south-via-relation",
   "Name": "rooms-south-via-relation",
   "Style": "ManyToOne",
   "Source": "rooms-south-via",
   "Target": "doors-x-via-south"
  },
  "rooms-up-via-relation": {
   "Id": "rooms-up-via-relation",
   "Name": "rooms-up-via-relation",
   "Style": "ManyToOne",
   "Source": "rooms-up-via",
   "Target": "doors-x-via-up"
  },
  "rooms-west-via-relation": {
   "Id": "rooms-west-via-relation",
   "Name": "rooms-west-via-relation",
   "Style": "ManyToOne",
   "Source": "rooms-west-via",
   "Target": "doors-x-via-west"
  },
  "supporters-contents-relation": {
   "Id": "supporters-contents-relation",
   "Name": "supporters-contents-relation",
   "Style": "OneToMany",
   "Source": "supporters-contents",
   "Target": "objects-support"
  }
 },
 "SingleToPlural": {
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
