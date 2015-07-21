package stories

import (
	. "github.com/ionous/sashimi/extensions"
	G "github.com/ionous/sashimi/game"
	. "github.com/ionous/sashimi/script"
)

//"First Conversation for the Alien Boy."
func Boy_Talk(s *Script) {
	s.The("room",
		Called("The Observatory"),
		Has("description", "An Empty Room"))

	s.The("actor", Called("The Alien Boy"), Exists())
	s.The("alien boy", Has("greeting", "OldBoy: What's The Matter?"))

	s.The("facts",
		Table("name", "summary").Has(
			"envious-phone", "I'm not too young for a phone, am I?").And(
			"boy-cam", "Camera? I hardly knew him.").And(
			"aliens-take-alot-of-photos", "Aliens sure take a lot of photos!").And(
			"camera-time", "My parents have allowed the impossible."))

	//Rule for beat-producing when the current interlocutor is the Alien Boy:
	//say `[one of]The Alien Boy scratches what seems to be his nose[or]The boy twirls a tentacle[or]The Alien Boy shifts on his stalks[or]The Alien Boy yawns first with one mouth, and then the other[or]The boy puckers a proboscis[at random].[run paragraph on]`.

	s.The("quip", Called("OldBoy: What's The Matter?"),
		Has("subject", "alien boy"),
		// Is("beat-opened"),
		Has("reply", `"You wouldn't happen to have a matter disrupter?", the Alien Boy asks.`),
		After("reporting discuss").Always(func(g G.Play) {
			g.The("alien boy").Set("greeting", g.The("OldBoy: Hello Again!"))
		}))

	s.The("quip", Called("OldBoy: Hello Again!"),
		Has("subject", "alien boy"),
		Is("repeatable"),
		Has("reply", `"Hello[one of] again[or][then at random]," the Alien Boy says.`))

	s.The("quip", Called("OldBoy: Later"),
		Has("subject", "alien boy"),
		Is("repeatable"),
		//Is("unimportant") -- can we determine that from the lack of reply?
		// or, even, it's position in the order of declaration?
		//Has("hook", `I"ll be back.`)
		Has("comment", `"Oh, sorry," Alice says. "I'll be back."`),
		// FIFIFIFIFIFIFIX! if this is "discuss" compilation succeeds, even though that's an impossible event.
		After("reporting discuss").Always(func(g G.Play) {
			g.The("player").Go("depart")
		}))

	s.The("quip requirements",
		Table("fact", "permitted-property", "quip").
			Has("boy-cam", "prohibited", "OldBoy: Later").
			And("boy-cam", "permitted", "OldBoy: CameraTime"))

	s.The("quip", Called("OldBoy: Camera"),
		Has("subject", "Alien Boy"),
		//Has("hook", `[if immediately]Oh really, ...[else]They made you get the camera?[end if]`)
		DirectlyFollows("OldBoy: Photos"),
		Has("comment", `"Wait," says Alice. "Your parents made you go get the camera?" says Alice.`),
		Has("reply", `"Amazing, huh?" says the Alien Boy. "You wouldn"t believe--"`),
		After("reporting discuss").Always(func(g G.Play) {
			con := g.Global("conversation").(*Conversation)
			con.Memory.Learn(g.The("boy-cam"))
		}))

	s.The("quip", Called("OldBoy: CameraTime1"),
		Has("subject", "Alien Boy"),
		DirectlyFollows("OldBoy: Camera"),
		//Has("hook",  `From your cabin?`)
		Has("reply", `"... Yes ..."`))

	s.The("quip", Called("OldBoy: CameraTime2"),
		Has("subject", "Alien Boy"),
		DirectlyFollows("OldBoy: CameraTime1"),
		//Has("hook",  `On your own?`),
		Has("reply", `"... Yes ..."`))

	s.The("quip", Called("OldBoy: CameraTime3"), Has("subject", "Alien Boy"),
		DirectlyFollows("OldBoy: CameraTime2"),
		//Has("hook",  `Where there are.... Toons?!!?`)
		Has("reply", `"Sometimes?" The Alien Boy looks confused.`))

	s.The("quip", Called("OldBoy: CameraTime"),
		Has("subject", "Alien Boy"),
		Is("repeatable"),
		//Is("unimportant"),
		//Is("weakly-phrased"),
		//Has("hook",  `Gotta run.`) ,
		Has("comment", `"Excuse me,"Alice says. "I need to see some parents about a cabin. Err... camera."`),
		Has("reply", `[beat] "I feel like there's some context I"m missing here..."`),
		After("reporting discuss").Always(func(g G.Play) {
			g.The("actor").Go("depart")
		}))

	s.The("quip",
		Called("OldBoy: DoesAnybody"),
		Has("subject", "Alien Boy"),
		DirectlyFollows("OldBoy: What's The Matter?"),
		//Is("super important"),  FIX: automatically sort "directly follows to the top.
		//Has("hook",  `A matter disrupter?`)
		Has("comment", `"A matter disrupter?" asks Alice. "Does anybody?"`),
		Has("reply", `"Or,"asks the Alien Boy, "maybe a ray gun?"`))

	s.The("quip",
		Called("OldBoy: Forcing"),
		Has("subject", "Alien Boy"),
		IndirectlyFollows("OldBoy: What's The Matter?"),
		IndirectlyFollows("OldBoy: Hello Again!"),
		//Has("hook",  `Watching the storm?`)
		Has("comment", `"[if immediately]Oh, are[else]Are[end if] your parents forcing you to watch the storm, too?"Alice asks.`),
		Has("reply", `"It"s not so bad. [beat] At least they didn"t lock me in the cabin like the last time."`))

	s.The("quip",
		Called("OldBoy: reverse-reply"),
		Has("subject", "AlienBoy"),
		Has("reply", `Another camera flash goes off, and Alice says "Your parents sure do take a lot of photos."`))

	s.The("quip",
		Called("OldBoy: Photos"),
		Has("subject", "Alien Boy"),
		//Has("hook",  `That"s a lot of photos!`),
		Has("comment", `"Your parents sure take a lot of photos,"Alice whispers conspiratorially.`),
		Has("reply", `"They"re nuts for it," says the Alien Boy. "Once my mom left hers and made [bold type]me[roman type] go get it."`))

	s.The("quip",
		Called("OldBoy: Cabin"),
		Has("subject", "Alien Boy"),
		IndirectlyFollows("OldBoy: Forcing"),
		//Has("hook",  `You got locked in the cabin?`),
		Has("comment", `"[if immediately]They locked [else]Did your parents really lock [end if]you in your cabin?" asks Alice incredulously. "Tell me you at least had Toons!?"`),
		Has("reply", `"Well..."[beat] "I [bold type]was[roman type] molting...."`))
	/*
	   alien-photo-turn is a number that varies"),
	   alien-photo-count is a number that varies"),

	   To say obs-random-photograph:
	   	if the current interlocutor is the Alien Boy:
	   		say `The boy's parents snap a photo [one of]of you[or]of him[or]of a passing asteroid[purely at random].`;
	   	otherwise:
	   		say `[one of]The alien family near the windows take another photo[or]Another camera flash goes off[or]You hear the sound of an ancient camera shutter[purely at random].`

	   To say obs-notices-photo:
	   	if the current interlocutor is the Alien Boy:
	   		queue OldBoy: reverse-reply as postponed obligatory;
	   	otherwise:
	   		say `Another camera flash goes off, and Alice says "That boy's parents sure do take a lot of photos.`;
	   	now the player knows aliens-take-alot-of-photos"),

	   Every turn when the player is in the observatory and the ship status is full steam ahead and the alien-photo-turn is less than the turn count:
	   	say `[one of]You see the flash of a camera out of the corner of your eye.[or][obs-random-photograph][or][obs-notices-photo][or][obs-random-photograph][stopping]`;
	   	if the player knows aliens-take-alot-of-photos:
	   		now alien-photo-turn is turn count plus a random number from two to four;
	   	otherwise:
	   		now alien-photo-turn is turn count plus a random number from zero to two"),
	*/

	s.The("quip requirements",
		Table("fact", "permitted-property", "quip").
			Has("aliens-take-alot-of-photos", "permitted", "OldBoy: Photos"))

	s.The("quip",
		Called("OldBoy: Molting"),
		Has("subject", "Alien Boy"),
		IndirectlyFollows("OldBoy: Cabin"),
		//Has("hook",  `Ik. You molt?` ),
		Has("comment", `"I probably don"t want to ask," says Alice, "do I?"`),
		Has("reply", `[beat] "Probably not. No."`))

	s.The("quip",
		Called("OldBoy: RayGun"),
		Has("subject", "Alien Boy"),
		DirectlyFollows("OldBoy: DoesAnybody"),
		//Is("really important"),
		//Has("hook",  `Ummm"... a ray gun?`),
		Has("comment", `"And what"re you gonna do with a ray gun?" Alice asks.`),
		Has("reply", `"I was thinking of using to get out of here."[beat]`))

	s.The("quip",
		Called("OldBoy: SeeNoEvil"),
		Has("subject", "Alien Boy"),
		DirectlyFollows("OldBoy: RayGun"),
		//Has("hook",  `Hear no evil, see no evil.`),
		Has("comment", `Alice says, "I probably don"t want to ask, do I?"`),
		Has("reply", `"My dad says violence never solves anything." The Alien Boy's antennae start shaking. "I"ll show him,"he yells. "I"ll show them all!"`))

	s.The("quip",
		Called("OldBoy: Yelp"),
		Has("subject", "Alien Boy"),
		DirectlyFollows("OldBoy: SeeNoEvil"),
		//Has("hook",  `Safety in distance.`),
		Has("comment", `"Ummm....", Alice says. "I"ll just be over here then."`),
		Has("reply", `[beat]"Was it something I said?"`),
		After("reporting discuss").Always(func(g G.Play) {
			g.The("actor").Go("depart")
		}))
}

func boyTalkStory(s *Script) {
	s.The("story",
		Called("boy talk"),
		Has("author", "me"))
	s.The("alien boy", Exists(), In("the observatory"))
	Boy_Talk(s)
}

func init() {
	stories.Register("boy", boyTalkStory)
}
