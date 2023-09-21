package core

import "github.com/guojia99/my-cubing/src/core/model"

type (
	Core interface {
		Read
		ReadPlayer
		ReadContest

		// ReloadCache 重置缓存
		ReloadCache()
		// AddScore 添加成绩
		AddScore(playerID uint, contestID uint, project model.Project, roundId uint, result []float64, penalty model.ScorePenalty) error
		// RemoveScore 删除成绩
		RemoveScore(scoreID uint) error
		// StatisticalRecordsAndEndContest 结束比赛并统计记录
		StatisticalRecordsAndEndContest(contestId uint) error
	}

	ReadPlayer interface {
		// GetPlayerDetail 获取玩家参赛信息
		GetPlayerDetail(playerId uint) PlayerDetail
		// GetPlayerBestScore 获取玩家最佳成绩及排名
		GetPlayerBestScore(playerId uint) (bestSingle, bestAvg map[model.Project]RankScore)
		// GetAllPlayerBestScore 获取所有人最佳成绩
		GetAllPlayerBestScore() (bestSingle, bestAvg map[model.Project][]model.Score)
		// GetAllPlayerBestScoreByProject 获取某个项目最佳成绩
		GetAllPlayerBestScoreByProject(project model.Project) (bestSingle, bestAvg []model.Score)
		// GetPodiumsByPlayer 获取玩家领奖台数据
		GetPodiumsByPlayer(playerID uint) Podiums
		// GetRecordByPlayer 获取一个人的记录
		GetRecordByPlayer(playerID uint) []RecordMessage
		// GetPlayerScore 获取选手所有成绩
		GetPlayerScore(playerID uint) (bestSingle, bestAvg []model.Score, scores []ScoresByContest)
		// GetPlayerSor 获取选手Sor数据
		GetPlayerSor(playerID uint) (single, avg map[model.SorStatisticsKey]SorScore)
	}

	ReadContest interface {
		// GetSorScoreByContest 获取某场比赛的排名总和
		GetSorScoreByContest(contestID uint) (single, avg map[model.SorStatisticsKey][]SorScore)
		// GetScoreByContest 获取某场比赛成绩排名
		GetScoreByContest(contestID uint) map[model.Project][]RoutesScores
		// GetPodiumsByContest 获取比赛的领奖台数据
		GetPodiumsByContest(contestID uint) []Podiums
	}

	Read interface {
		// GetBestScores 获取所有项目最佳成绩
		GetBestScores() (bestSingle, bestAvg map[model.Project]model.Score)
		// GetSorScore 获取排名总和
		GetSorScore() (single, avg map[model.SorStatisticsKey][]SorScore)
		// GetAllPodium 获取全部人的领奖台排行
		GetAllPodium() []Podiums
		// GetRecordByContest 获取一场比赛中的记录
		GetRecordByContest(contestID uint) []RecordMessage
	}
)
