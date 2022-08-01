package main

//斗地主

type CardsType = byte

const (
	error       CardsType = iota //错误
	single                       //单张
	double                       //一对
	three                        //飞机
	four                         //炸弹
	twoKing                      //王炸
	threeSingle                  //三带一(或六带二，九带三。。。。)
	threeDouble                  //三带一对(或六带四，九带六。。。。)
	fourSingle                   //四带二(或八带四，十二带六。。。。)
	fourDouble                   //四带两对(或八带四对，十二带六。。。。)
	singles                      //大于5张的连牌
	doubles                      //大于三对的对对连牌
	threes                       //大于两个飞机的飞机连牌
	foures                       //双连炸
	foures1                      //三连炸
	foures2                      //四连炸
	foures3                      //五连炸
)

// Analyse 分析牌型
type analyseCard struct {
	CardsType [4]byte
	Cards     [4][20]byte
}

type FightLandlords struct {
}

func getCardColor(card byte) byte {
	return card & 0xF0
}
func getCardValue(card byte) byte {
	return card & 0x0F
}

func getCardLogicValue(card byte) byte {
	cardColor := getCardColor(card)
	cardValue := getCardValue(card)

	if cardColor == 0x40 {
		return cardValue + 2
	}

	if cardValue <= 2 {
		return cardValue + 13
	}
	return cardValue
}

func (analyse *analyseCard) Analyse(cards []byte) {
	sameCount := 0
	logicValue := byte(0)
	cardsLen := len(cards)
	for i := 0; i < cardsLen; i++ {
		sameCount = 1
		logicValue = getCardLogicValue(cards[i])
		for j := i + 1; j < cardsLen; j++ {
			if getCardLogicValue(cards[j]) != logicValue {
				break
			}
			sameCount++
		}
		analyse.CardsType[sameCount-1]++

		for j := 0; j < sameCount; j++ {
			analyse.Cards[sameCount-1][analyse.CardsType[sameCount-1]*byte(sameCount)+byte(j)] = cards[i+j]
		}
		i = i + sameCount - 1
	}
}

// Sort 按大小排序
func (fightlandlords *FightLandlords) Sort(cards []byte) {
	cardsCount := len(cards)
	for i := 0; i < cardsCount; i++ {
		for j := i + 1; j < cardsCount; j++ {
			if getCardLogicValue(cards[i]) > getCardLogicValue(cards[j]) {
				cards[i], cards[j] = cards[j], cards[i]
			}
		}
	}
}

// CheckCardsType 获取牌型
func (fightlandlords *FightLandlords) CheckCardsType(cards []byte) CardsType {
	cardsCount := byte(len(cards))
	switch cardsCount {
	case 1:
		return single
	case 2:
		if (cards[0] == 0x4F) && (cards[1] == 0x4E) {
			return twoKing //如果是大王和小王，返回火箭
		}

		if getCardLogicValue(cards[0]) == getCardLogicValue(cards[1]) {
			return double //如果两张牌逻辑值相同，返回一队
		}
	}
	analyse := &analyseCard{}
	analyse.Analyse(cards)
	//带有四张炸弹的牌判断
	if analyse.CardsType[3] > 0 {
		if (analyse.CardsType[3] == 1) && (cardsCount == 4) {
			return four
		}
		if (analyse.CardsType[3] == 1) && (cardsCount == 6) {
			return fourSingle
		}
		if (analyse.CardsType[3] == 1) && (cardsCount == 8) && (analyse.CardsType[1] == 2) {
			return fourDouble
		}
		if (analyse.CardsType[3] == 2) && (cardsCount == 8) {
			return foures
		}
		if (analyse.CardsType[3] == 3) && (cardsCount == 12) {
			return foures1
		}
		if (analyse.CardsType[3] == 4) && (cardsCount == 16) {
			return foures2
		}
		if (analyse.CardsType[3] == 5) && (cardsCount == 20) {
			return foures3
		}
	}
	//如果飞机数大于0
	if analyse.CardsType[2] > 0 {
		if analyse.CardsType[2] > 1 { //如果大于1个飞机
			firstValue := getCardLogicValue(analyse.Cards[2][0])
			if firstValue >= 15 {
				return error
			}
			//连牌判断
			for i := byte(1); i < analyse.CardsType[2]; i++ {
				if firstValue != (getCardLogicValue(analyse.Cards[2][i*3]) + i) {
					return error
				}
			}
		} else if cardsCount == 3 {
			return three
		}
		if analyse.CardsType[2]*3 == cardsCount {
			return threes
		}
		if analyse.CardsType[2]*4 == cardsCount {
			return threeSingle
		}
		if (analyse.CardsType[2]*5 == cardsCount) && (analyse.CardsType[1] == analyse.CardsType[2]) {
			return threeDouble
		}

	}
	//如果对子数大于3
	if analyse.CardsType[1] >= 3 {
		firstValue := getCardLogicValue(analyse.Cards[1][0])
		if firstValue >= 15 {
			return error
		}
		for i := byte(1); i < analyse.CardsType[1]; i++ {
			if firstValue != (getCardLogicValue(analyse.Cards[1][i*2]) + i) {
				return error
			}
		}
		if (analyse.CardsType[1] * 2) == cardsCount {
			return doubles
		}
	}
	//单连判断
	if (analyse.CardsType[0] >= 5) && (analyse.CardsType[0] == cardsCount) {
		firstValue := getCardLogicValue(analyse.Cards[0][0])
		if firstValue >= 15 {
			return error
		}
		//连牌判断
		for i := byte(1); i < analyse.CardsType[0]; i++ {
			if firstValue != (getCardLogicValue(analyse.Cards[0][i]) + i) {
				return error
			}
		}
		return singles
	}
	return error
}

// Compare 比较两个牌大小
func (fightlandlords *FightLandlords) Compare(cards1, cards2 []byte) {

}
