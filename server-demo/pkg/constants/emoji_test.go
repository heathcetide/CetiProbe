package constants

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmojiConstants(t *testing.T) {
	// 测试表情符号常量不为空
	assert.NotEmpty(t, Happy)
	assert.NotEmpty(t, Smile)
	assert.NotEmpty(t, LoveEye)
	assert.NotEmpty(t, Cry)
	assert.NotEmpty(t, Neutral)
	assert.NotEmpty(t, Confused)
	assert.NotEmpty(t, Worried)
	assert.NotEmpty(t, Angry)
	assert.NotEmpty(t, Rage)
	assert.NotEmpty(t, Sweat)
	assert.NotEmpty(t, Tired)
	assert.NotEmpty(t, Sleepy)
	assert.NotEmpty(t, Unamused)
	assert.NotEmpty(t, RollingEyes)
	assert.NotEmpty(t, Blush)
	assert.NotEmpty(t, Grin)
	assert.NotEmpty(t, Joy)
	assert.NotEmpty(t, SlightSmile)
	assert.NotEmpty(t, Wink)
	assert.NotEmpty(t, Kiss)
	assert.NotEmpty(t, StarStruck)
	assert.NotEmpty(t, Thinking)
	assert.NotEmpty(t, Yawning)
}

func TestHandEmojiConstants(t *testing.T) {
	// 测试手势表情符号
	assert.NotEmpty(t, ThumbsUp)
	assert.NotEmpty(t, ThumbsDown)
	assert.NotEmpty(t, Victory)
	assert.NotEmpty(t, CallMe)
	assert.NotEmpty(t, Wave)
	assert.NotEmpty(t, FoldedHands)
	assert.NotEmpty(t, Clap)
	assert.NotEmpty(t, RaisedHand)
	assert.NotEmpty(t, OkHand)
}

func TestWeatherEmojiConstants(t *testing.T) {
	// 测试天气相关表情符号
	assert.NotEmpty(t, Sun)
	assert.NotEmpty(t, Moon)
	assert.NotEmpty(t, Cloud)
	assert.NotEmpty(t, Rain)
	assert.NotEmpty(t, Lightning)
	assert.NotEmpty(t, Snowflake)
	assert.NotEmpty(t, Fire)
}

func TestNatureEmojiConstants(t *testing.T) {
	// 测试自然相关表情符号
	assert.NotEmpty(t, Tree)
	assert.NotEmpty(t, Flower)
}

func TestAnimalEmojiConstants(t *testing.T) {
	// 测试动物表情符号
	assert.NotEmpty(t, Dog)
	assert.NotEmpty(t, Cat)
	assert.NotEmpty(t, Monkey)
	assert.NotEmpty(t, Lion)
	assert.NotEmpty(t, Tiger)
	assert.NotEmpty(t, Elephant)
	assert.NotEmpty(t, Panda)
	assert.NotEmpty(t, Horse)
	assert.NotEmpty(t, Cow)
	assert.NotEmpty(t, Fish)
}

func TestFoodEmojiConstants(t *testing.T) {
	// 测试食物表情符号
	assert.NotEmpty(t, Pizza)
	assert.NotEmpty(t, Burger)
	assert.NotEmpty(t, Sushi)
	assert.NotEmpty(t, Taco)
	assert.NotEmpty(t, Hotdog)
	assert.NotEmpty(t, Cake)
	assert.NotEmpty(t, IceCream)
	assert.NotEmpty(t, Coffee)
	assert.NotEmpty(t, Beer)
	assert.NotEmpty(t, WineGlass)
}

func TestTechnologyEmojiConstants(t *testing.T) {
	// 测试科技相关表情符号
	assert.NotEmpty(t, Computer)
	assert.NotEmpty(t, Mobile)
	assert.NotEmpty(t, Camera)
	assert.NotEmpty(t, Headphones)
}

func TestToolEmojiConstants(t *testing.T) {
	// 测试工具相关表情符号
	assert.NotEmpty(t, Tools)
	assert.NotEmpty(t, Hammer)
	assert.NotEmpty(t, Wrench)
	assert.NotEmpty(t, Gear)
	assert.NotEmpty(t, Microscope)
	assert.NotEmpty(t, Rocket)
}

func TestSymbolEmojiConstants(t *testing.T) {
	// 测试符号表情符号
	assert.NotEmpty(t, Checkmark)
	assert.NotEmpty(t, CrossMark)
	assert.NotEmpty(t, Info)
	assert.NotEmpty(t, Warning)
	assert.NotEmpty(t, NoEntry)
	assert.NotEmpty(t, Bell)
	assert.NotEmpty(t, Lock)
	assert.NotEmpty(t, Unlock)
}

func TestEmojiUnicodeFormat(t *testing.T) {
	// 测试表情符号的Unicode格式
	// 大部分表情符号应该以\U开头（32位Unicode）
	unicodeEmojis := []string{
		Happy, Smile, LoveEye, Neutral, Confused, Worried,
		Angry, Rage, Sweat, Tired, Sleepy, Unamused, RollingEyes, Blush,
		Grin, Joy, SlightSmile, Wink, Kiss, StarStruck, Thinking,
		Yawning, ThumbsUp, ThumbsDown, Victory, CallMe, Wave, FoldedHands,
		Clap, RaisedHand, OkHand, Sun, Moon, Cloud, Rain, Lightning,
		Snowflake, Fire, Tree, Flower, Dog, Cat, Monkey, Lion, Tiger,
		Elephant, Panda, Horse, Cow, Fish, Pizza, Burger, Sushi, Taco,
		Hotdog, Cake, IceCream, Coffee, Beer, WineGlass, Computer, Mobile,
		Camera, Headphones, Tools, Hammer, Wrench, Gear, Microscope, Rocket,
		Checkmark, CrossMark, Info, Warning, NoEntry, Bell, Lock, Unlock,
	}

	for _, emoji := range unicodeEmojis {
		// 检查是否以\U开头（32位Unicode）或\u开头（16位Unicode）
		assert.True(t, len(emoji) >= 6, "表情符号 %s 长度应该至少为6", emoji)
		assert.True(t, emoji[0] == '\\', "表情符号 %s 应该以\\开头", emoji)
		assert.True(t, emoji[1] == 'U' || emoji[1] == 'u', "表情符号 %s 应该以\\U或\\u开头", emoji)
	}
}

func TestEmojiUniqueness(t *testing.T) {
	// 测试表情符号常量的唯一性
	allEmojis := []string{
		Happy, Smile, LoveEye, Cry, Neutral, Confused,
		Worried, Angry, Rage, Sweat, Tired, Sleepy, Unamused, RollingEyes,
		Blush, Grin, Joy, SlightSmile, Wink, Kiss, StarStruck,
		Thinking, Yawning, ThumbsUp, ThumbsDown, Victory, CallMe, Wave,
		FoldedHands, Clap, RaisedHand, OkHand, Sun, Moon, Cloud, Rain,
		Lightning, Snowflake, Fire, Tree, Flower, Dog, Cat, Monkey, Lion,
		Tiger, Elephant, Panda, Horse, Cow, Fish, Pizza, Burger, Sushi,
		Taco, Hotdog, Cake, IceCream, Coffee, Beer, WineGlass, Computer,
		Mobile, Camera, Headphones, Tools, Hammer, Wrench, Gear, Microscope,
		Rocket, Checkmark, CrossMark, Info, Warning, NoEntry, Bell, Lock, Unlock,
	}

	// 检查是否有重复的表情符号
	seen := make(map[string]bool)
	for _, emoji := range allEmojis {
		assert.False(t, seen[emoji], "表情符号 %s 重复了", emoji)
		seen[emoji] = true
	}
}

func TestSpecialEmojiCases(t *testing.T) {
	// 测试特殊情况
	// Cry表情符号使用不同的格式
	assert.Equal(t, "\uF622", Cry)

	// 验证某些表情符号的特定值
	assert.Equal(t, "\\U0001F603", Happy)
	assert.Equal(t, "\\U0001F604", Smile)
	assert.Equal(t, "\\U0001F60D", LoveEye)
}

func TestEmojiCategories(t *testing.T) {
	// 测试表情符号分类
	faceEmojis := []string{Happy, Smile, LoveEye, Cry, Neutral, Confused, Worried, Angry, Rage, Sweat, Tired, Sleepy, Unamused, RollingEyes, Blush, Grin, Joy, SlightSmile, Wink, Kiss, StarStruck, Thinking, Yawning}
	handEmojis := []string{ThumbsUp, ThumbsDown, Victory, CallMe, Wave, FoldedHands, Clap, RaisedHand, OkHand}
	weatherEmojis := []string{Sun, Moon, Cloud, Rain, Lightning, Snowflake, Fire}
	animalEmojis := []string{Dog, Cat, Monkey, Lion, Tiger, Elephant, Panda, Horse, Cow, Fish}
	foodEmojis := []string{Pizza, Burger, Sushi, Taco, Hotdog, Cake, IceCream, Coffee, Beer, WineGlass}
	symbolEmojis := []string{Checkmark, CrossMark, Info, Warning, NoEntry, Bell, Lock, Unlock}

	// 验证每个分类都不为空
	assert.NotEmpty(t, faceEmojis)
	assert.NotEmpty(t, handEmojis)
	assert.NotEmpty(t, weatherEmojis)
	assert.NotEmpty(t, animalEmojis)
	assert.NotEmpty(t, foodEmojis)
	assert.NotEmpty(t, symbolEmojis)
}
